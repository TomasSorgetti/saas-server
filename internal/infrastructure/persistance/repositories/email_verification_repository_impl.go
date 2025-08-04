package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"luthierSaas/internal/domain/entities"
	"luthierSaas/internal/interfaces/repository"
)

type emailVerificationRepository struct {
	db *sql.DB
}

func NewEmailVerificationRepository(db *sql.DB) repository.EmailVerificationRepository {
	return &emailVerificationRepository{db: db}
}

func (r *emailVerificationRepository) Create(ctx context.Context, ev *entities.EmailVerification) error {
	query := `INSERT INTO email_verifications (user_id, verification_code, expires_at, verified)
			  VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, ev.UserID, ev.Code, ev.ExpiresAt, ev.Verified)
	return err
}

func (r *emailVerificationRepository) GetByUserID(ctx context.Context, userID int) (*entities.EmailVerification, error) {
	query := `SELECT id, user_id, verification_code, expires_at, verified, created_at, updated_at
			  FROM email_verifications WHERE user_id = ?`
	row := r.db.QueryRowContext(ctx, query, userID)

	var ev entities.EmailVerification
	err := row.Scan(
		&ev.ID,
		&ev.UserID,
		&ev.Code,
		&ev.ExpiresAt,
		&ev.Verified,
		&ev.CreatedAt,
		&ev.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return &ev, err
}

func (r *emailVerificationRepository) MarkAsVerified(ctx context.Context, userID int) error {
	query := `UPDATE email_verifications SET verified = TRUE WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *emailVerificationRepository) DeleteByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM email_verifications WHERE user_id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}

func (r *emailVerificationRepository) UpdateCode(ctx context.Context, id int, newCode string, newExpiresAt time.Time) error {
	query := `UPDATE email_verifications SET verification_code = ?, expires_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, newCode,newExpiresAt, id)
	return err
}