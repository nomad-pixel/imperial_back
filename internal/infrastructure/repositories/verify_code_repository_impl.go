package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type VerifyCodeRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewVerifyCodeRepositoryImpl(db *pgxpool.Pool) ports.VerifyCodeRepository {
	return &VerifyCodeRepositoryImpl{db: db}
}

func (r *VerifyCodeRepositoryImpl) CreateVerifyCode(ctx context.Context, code string, userID int64, verifyCodeType entities.VerifyCodeType, expiresAt time.Time) (*entities.VerifyCode, error) {
	query := `
	INSERT INTO verify_codes (code, user_id, type, expires_at)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT ON CONSTRAINT verify_codes_user_id_type_key DO UPDATE
	SET
			code       = EXCLUDED.code,
			expires_at = EXCLUDED.expires_at,
			is_used    = FALSE,
			updated_at = NOW()
	RETURNING id, code, user_id, type, is_used, expires_at, created_at, updated_at;
	`
	var verifyCode entities.VerifyCode
	err := r.db.QueryRow(ctx, query, code, userID, verifyCodeType, expiresAt).Scan(&verifyCode.ID, &verifyCode.Code, &verifyCode.UserID, &verifyCode.Type, &verifyCode.IsUsed, &verifyCode.ExpiresAt, &verifyCode.CreatedAt, &verifyCode.UpdatedAt)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &verifyCode, nil
}

func (r *VerifyCodeRepositoryImpl) GetVerifyCodeByUserIDAndType(ctx context.Context, userID int64, verifyCodeType entities.VerifyCodeType) (*entities.VerifyCode, error) {
	query := `
		SELECT id, code, user_id, type, is_used, expires_at, created_at, updated_at
		FROM verify_codes
		WHERE user_id = $1 AND type = $2
	`
	var verifyCode entities.VerifyCode
	err := r.db.QueryRow(ctx, query, userID, verifyCodeType).Scan(&verifyCode.ID, &verifyCode.Code, &verifyCode.UserID, &verifyCode.Type, &verifyCode.IsUsed, &verifyCode.ExpiresAt, &verifyCode.CreatedAt, &verifyCode.UpdatedAt)
	if err != nil {
		return nil, r.handleError(err)
	}
	return &verifyCode, nil
}

func (r *VerifyCodeRepositoryImpl) UpdateVerifyCode(ctx context.Context, verifyCode *entities.VerifyCode) (*entities.VerifyCode, error) {
	query := `
		UPDATE verify_codes
		SET is_used = $1, expires_at = $2, code = $3
		WHERE id = $4
	`
	err := r.db.QueryRow(ctx, query, verifyCode.IsUsed, verifyCode.ExpiresAt, verifyCode.Code, verifyCode.ID).Scan(&verifyCode.ID, &verifyCode.Code, &verifyCode.UserID, &verifyCode.Type, &verifyCode.IsUsed, &verifyCode.ExpiresAt, &verifyCode.CreatedAt, &verifyCode.UpdatedAt)
	if err != nil {
		return nil, r.handleError(err)
	}
	return verifyCode, nil
}

func (r *VerifyCodeRepositoryImpl) handleError(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return apperrors.ErrVerifyCodeNotFound
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return apperrors.Wrap(err, apperrors.ErrCodeConflict, "Код верификации уже существует")
		case "23503": // foreign_key_violation
			return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Нарушение внешнего ключа")
		case "23514": // check_violation
			return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Нарушение ограничения проверки")
		}
	}

	return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Ошибка при работе с базой данных")
}
