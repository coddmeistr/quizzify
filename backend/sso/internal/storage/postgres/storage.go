package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/domain/models"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/storage"
)

type Storage struct {
	db *pgx.Conn
}

func New(connectionString string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, login string, email string, passHash []byte) (uint64, error) {
	const op = "storage.postgres.SaveUser"

	var lastInsertedId uint64
	if err := s.db.QueryRow(ctx, "INSERT INTO users(login, email, pass_hash) VALUES($1, $2, $3) RETURNING id", login, email, passHash).Scan(&lastInsertedId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // // Unique violation
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return lastInsertedId, nil
}

func (s *Storage) UserByLogin(ctx context.Context, login string) (models.User, error) {
	const op = "storage.postgres.UserByLogin"

	user := models.User{}
	if err := s.db.QueryRow(ctx, "SELECT id, login, email, pass_hash FROM users WHERE login = $1", login).Scan(&user.ID, &user.Login, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.UserByEmail"

	user := models.User{}
	if err := s.db.QueryRow(ctx, "SELECT id, login, email, pass_hash FROM users WHERE email = $1", email).Scan(&user.ID, &user.Login, &user.Email, &user.PassHash); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID uint64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	var isAdmin bool
	if err := s.db.QueryRow(ctx, "SELECT is_admin FROM users WHERE id = $1", userID).Scan(&isAdmin); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.postgres.App"

	app := models.App{}
	if err := s.db.QueryRow(ctx, "SELECT id, name, secret FROM apps WHERE id = $1", appID).Scan(&app.ID, &app.Name, &app.Secret); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return app, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
