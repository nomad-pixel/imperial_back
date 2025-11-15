package repositories

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type UserRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepositoryImpl(db *pgxpool.Pool) ports.UserRepository {
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, email, password string) (*entities.User, error) {
	query := `
		INSERT INTO users (email, password_hash, verified_at)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, verified_at, created_at, updated_at
	`
	var user entities.User
	err := r.db.QueryRow(ctx, query, email, password, false).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.VerifiedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT id, email, password_hash, verified_at, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	var user entities.User
	err := r.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.VerifiedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUserById(ctx context.Context, id int64) (*entities.User, error) {
	query := `
		SELECT id, email, password_hash, verified_at, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	var user entities.User
	err := r.db.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.VerifiedAt, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, updated_at = NOW()
		WHERE id = $3
		RETURNING id, email, password_hash, verified_at, created_at, updated_at
	`
	var updatedUser entities.User
	err := r.db.QueryRow(ctx, query, user.Email, user.PasswordHash, user.ID).Scan(&updatedUser.ID, &updatedUser.Email, &updatedUser.PasswordHash, &updatedUser.VerifiedAt, &updatedUser.CreatedAt, &updatedUser.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &updatedUser, nil
}

func (r *UserRepositoryImpl) DeleteUser(ctx context.Context, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`
	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return r.handleError(err)
	}
	return nil
}

// handleError преобразует ошибки БД в доменные ошибки
func (r *UserRepositoryImpl) handleError(err error) error {
	// Обработка "не найдено"
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrUserNotFound
	}

	// Обработка PostgreSQL специфичных ошибок
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return apperrors.ErrUserAlreadyExists
		case "23503": // foreign_key_violation
			return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Нарушение внешнего ключа")
		case "23514": // check_violation
			return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Нарушение ограничения проверки")
		}
	}

	// Остальные ошибки оборачиваем как database errors
	return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Ошибка при работе с базой данных")
}
