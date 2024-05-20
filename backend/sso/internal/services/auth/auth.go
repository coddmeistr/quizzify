package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/coddmeistr/quizzify/backend/sso/internal/domain/models"
	appjwt "github.com/coddmeistr/quizzify/backend/sso/internal/lib/jwt"
	"github.com/coddmeistr/quizzify/backend/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserProvider interface {
	UserByID(ctx context.Context, id uint64) (models.User, error)
	UserByLogin(ctx context.Context, login string) (models.User, error)
	UserByEmail(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID uint64) (bool, error)
	DeleteUser(ctx context.Context, userID uint64) error
	ListUsers(ctx context.Context) ([]models.User, error)
}

type UserSaver interface {
	SaveUser(ctx context.Context, login string, email string, passHash []byte) (uint64, error)
}

type PermissionsProvider interface {
	UserPermissions(ctx context.Context, userID int64) ([]int, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserExists         = errors.New("user exists")
	ErrAppNotFound        = errors.New("app not found")
)

type Auth struct {
	log           *slog.Logger
	usrSaver      UserSaver
	usrProvider   UserProvider
	permsProvider PermissionsProvider
	appProvider   AppProvider
	tokenTTL      time.Duration
}

func New(
	log *slog.Logger,
	usrSaver UserSaver,
	usrProvider UserProvider,
	permsProvider PermissionsProvider,
	appProvider AppProvider,
	tokenTTL time.Duration) *Auth {
	return &Auth{
		log:           log,
		usrSaver:      usrSaver,
		usrProvider:   usrProvider,
		permsProvider: permsProvider,
		appProvider:   appProvider,
		tokenTTL:      tokenTTL,
	}
}

func (a *Auth) AccountsList(ctx context.Context) ([]models.User, error) {
	users, err := a.usrProvider.ListUsers(ctx)
	return users, err
}

func (a *Auth) DeleteAccount(ctx context.Context, userID uint64) error {
	const op = "auth.DeleteAccount"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("user_id", int(userID)),
	)
	log.Info("deleting user")

	err := a.usrProvider.DeleteUser(ctx, userID)
	if err != nil {
		log.Error("failed deleting user", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user deleted")
	return nil
}

func (a *Auth) UserInfo(ctx context.Context, userID uint64) (models.User, []int, error) {
	const op = "auth.UserInfo"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("user_id", int(userID)),
	)
	log.Info("getting user info")

	user, err := a.usrProvider.UserByID(ctx, userID)
	if err != nil {
		log.Error("failed getting user info", slog.String("error", err.Error()))
		return models.User{}, nil, err
	}

	perms, err := a.permsProvider.UserPermissions(ctx, int64(userID))
	if err != nil {
		log.Error("failed getting user permissions", slog.String("error", err.Error()))
		return models.User{}, nil, err
	}

	return user, perms, nil
}

func (a *Auth) Login(ctx context.Context, login string, email string, password string, appID int) (token string, err error) {
	const op = "auth.Login"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("attempting to login user")

	// If login is provided then using login
	// In case if login is empty then using email
	var user models.User
	if login != "" {
		user, err = a.usrProvider.UserByLogin(ctx, login)
	} else {
		user, err = a.usrProvider.UserByEmail(ctx, email)
	}
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("failed to get user", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", op, err)
		} else {
			log.Error("failed getting current app", slog.String("error", err.Error()))
			return "", fmt.Errorf("%s: %w", op, err)
		}
	}

	log.Info("user logged in successfully")

	perms, err := a.permsProvider.UserPermissions(ctx, int64(user.ID))
	if err != nil {
		log.Error("failed getting user permissions", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err = appjwt.NewToken(user, perms, app, a.tokenTTL)
	if err != nil {
		log.Error("failed generating new token", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("token was generated")

	return token, nil
}

func (a *Auth) Register(ctx context.Context, login string, email string, password string) (userID uint64, err error) {
	const op = "auth.Register"
	log := a.log.With(
		slog.String("op", op),
	)
	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed generating hash from password", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	userID, err = a.usrSaver.SaveUser(ctx, login, email, passHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Warn("user already exists", slog.String("error", err.Error()))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		log.Error("failed saving user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user registered")
	return userID, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID uint64) (bool, error) {
	const op = "auth.IsAdmin"
	log := a.log.With(
		slog.String("op", op),
		slog.Int("user_id", int(userID)),
	)
	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Warn("user not found", slog.String("error", err.Error()))
			return false, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		log.Error("failed checking if user is admin", slog.String("error", err.Error()))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checking completed", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
