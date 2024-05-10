package permissions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/coddmeistr/quizzify/backend/sso/internal/storage"
)

type PermissionManipulator interface {
	AddPermission(ctx context.Context, userID int64, permID int64) error
	RemovePermission(ctx context.Context, userID int64, permID int64) error
}

type Permissions struct {
	log     *slog.Logger
	storage PermissionManipulator
}

var (
	ErrPermissionExist    = errors.New("permission already exist")
	ErrPermissionNotExist = errors.New("permission not exist")
)

func New(log *slog.Logger, pm PermissionManipulator) *Permissions {
	return &Permissions{
		log:     log,
		storage: pm,
	}
}

func (p *Permissions) AddPermission(ctx context.Context, userID int64, permID int64) error {
	const op = "permissions.AddPermission"
	log := p.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
		slog.Int64("permission_id", permID),
	)

	log.Info("Adding permission to user")

	if err := p.storage.AddPermission(ctx, userID, permID); err != nil {
		if errors.Is(err, storage.ErrPermissionAlreadyExist) {
			log.Warn("permission exist", slog.String("error", err.Error()))
			return fmt.Errorf("%s: %w", op, ErrPermissionExist)
		}

		log.Error("failed adding permission", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Permission added")
	return nil
}

func (p *Permissions) RemovePermission(ctx context.Context, userID int64, permID int64) error {
	const op = "permissions.RemovePermission"
	log := p.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
		slog.Int64("permission_id", permID),
	)

	log.Info("Removing permission from user")

	if err := p.storage.RemovePermission(ctx, userID, permID); err != nil {
		if errors.Is(err, storage.ErrNoPermission) {
			log.Warn("permission not exist", slog.String("error", err.Error()))
			return fmt.Errorf("%s: %w", op, ErrPermissionNotExist)
		}

		log.Error("failed removing permission", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Permission removed")
	return nil
}
