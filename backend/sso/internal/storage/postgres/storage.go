package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/domain/models"
	"github.com/maxik12233/quizzify-online-tests/backend/sso/internal/storage"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connectionString string) (*Storage, error) {
	const op = "storage.postgres.New"

	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: pool}, nil
}

func (s *Storage) UserPermissions(ctx context.Context, userID int64) ([]int, error) {
	const op = "storage.postgres.UserPermissions"

	rows, err := s.db.Query(ctx, "SELECT permission_id FROM user_permissions WHERE user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	perms := make([]int, 0)
	permID := 0
	for {
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return nil, fmt.Errorf("%s: %w", op, err)
			}
			break
		}
		if err := rows.Scan(&permID); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		perms = append(perms, permID)
	}

	return perms, nil
}

func (s *Storage) AddPermission(ctx context.Context, userID int64, permID int64) error {
	const op = "storage.postgres.AddPermission"

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.QueryRow(ctx, "SELECT * FROM user_permissions WHERE user_id = $1 AND permission_id = $2", userID, permID).Scan(); err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("%s: %w", op, storage.ErrPermissionAlreadyExist)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if err := tx.QueryRow(ctx, "INSERT INTO user_permissions(user_id, permission_id) VALUES($1, $2)", userID, permID).Scan(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	tx.Commit(ctx)
	return nil
}

func (s *Storage) RemovePermission(ctx context.Context, userID int64, permID int64) error {
	const op = "storage.postgres.RemovePermission"

	var rowsDeleted int
	if err := s.db.QueryRow(ctx, "DELETE FROM user_permissions WHERE user_id = $1 AND permission_id = $2", userID, permID).Scan(&rowsDeleted); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsDeleted != 1 {
		return fmt.Errorf("%s: %w", op, storage.ErrNoPermission)
	}

	return nil
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
