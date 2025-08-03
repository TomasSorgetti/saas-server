package repositories

import (
	"context"
	"database/sql"
	"errors"

	"luthierSaas/internal/domain/entities"
	"luthierSaas/internal/interfaces/repository"
)

type sessionRepository struct {
    db *sql.DB
}

func NewSessionRepository(db *sql.DB) repository.SessionRepository {
    return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(ctx context.Context, session *entities.Session) error {
    query := `INSERT INTO sessions (user_id, access_token_hash, refresh_token_hash, expires_at, refresh_expires_at, is_valid, device_info)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
    _, err := r.db.ExecContext(ctx, query,
        session.UserID,
        session.AccessTokenHash,
        session.RefreshTokenHash,
        session.ExpiresAt,
        session.RefreshExpiresAt,
        session.IsValid,
        session.DeviceInfo,
    )
    return err
}

func (r *sessionRepository) FindByAccessTokenHash(ctx context.Context, accessTokenHash string) (*entities.Session, error) {
    query := `SELECT id, user_id, access_token_hash, refresh_token_hash, expires_at, refresh_expires_at, is_valid, device_info, created_at, updated_at
              FROM sessions WHERE access_token_hash = ? AND is_valid = TRUE AND expires_at > NOW()`
    row := r.db.QueryRowContext(ctx, query, accessTokenHash)

    var session entities.Session
    var refreshExpiresAt sql.NullTime
    var updatedAt sql.NullTime
    var deviceInfo sql.NullString

    err := row.Scan(
        &session.ID,
        &session.UserID,
        &session.AccessTokenHash,
        &session.RefreshTokenHash,
        &session.ExpiresAt,
        &refreshExpiresAt,
        &session.IsValid,
        &deviceInfo,
        &session.CreatedAt,
        &updatedAt,
    )

    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    if refreshExpiresAt.Valid {
        session.RefreshExpiresAt = refreshExpiresAt.Time
    }
    if updatedAt.Valid {
        session.UpdatedAt = updatedAt.Time
    }
    if deviceInfo.Valid {
        session.DeviceInfo = deviceInfo.String
    }

    return &session, nil
}

func (r *sessionRepository) FindByRefreshTokenHash(ctx context.Context, refreshTokenHash string) (*entities.Session, error) {
    query := `SELECT id, user_id, access_token_hash, refresh_token_hash, expires_at, refresh_expires_at, is_valid, device_info, created_at, updated_at
              FROM sessions WHERE refresh_token_hash = ? AND is_valid = TRUE AND refresh_expires_at > NOW()`
    row := r.db.QueryRowContext(ctx, query, refreshTokenHash)

    var session entities.Session
    var refreshExpiresAt sql.NullTime
    var updatedAt sql.NullTime
    var deviceInfo sql.NullString

    err := row.Scan(
        &session.ID,
        &session.UserID,
        &session.AccessTokenHash,
        &session.RefreshTokenHash,
        &session.ExpiresAt,
        &refreshExpiresAt,
        &session.IsValid,
        &deviceInfo,
        &session.CreatedAt,
        &updatedAt,
    )

    if errors.Is(err, sql.ErrNoRows) {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    if refreshExpiresAt.Valid {
        session.RefreshExpiresAt = refreshExpiresAt.Time
    }
    if updatedAt.Valid {
        session.UpdatedAt = updatedAt.Time
    }
    if deviceInfo.Valid {
        session.DeviceInfo = deviceInfo.String
    }

    return &session, nil
}

func (r *sessionRepository) Update(ctx context.Context, session *entities.Session) error {
    query := `UPDATE sessions SET access_token_hash = ?, refresh_token_hash = ?, expires_at = ?, refresh_expires_at = ?, is_valid = ?, device_info = ?, updated_at = NOW()
              WHERE id = ?`
    _, err := r.db.ExecContext(ctx, query,
        session.AccessTokenHash,
        session.RefreshTokenHash,
        session.ExpiresAt,
        session.RefreshExpiresAt,
        session.IsValid,
        session.DeviceInfo,
        session.ID,
    )
    return err
}

func (r *sessionRepository) Invalidate(ctx context.Context, accessTokenHash string) error {
    query := `UPDATE sessions SET is_valid = FALSE, updated_at = NOW() WHERE access_token_hash = ?`
    _, err := r.db.ExecContext(ctx, query, accessTokenHash)
    return err
}

func (r *sessionRepository) FindByUserID(ctx context.Context, userID int64) ([]*entities.Session, error) {
    query := `SELECT id, user_id, access_token_hash, refresh_token_hash, expires_at, refresh_expires_at, is_valid, device_info, created_at, updated_at
              FROM sessions WHERE user_id = ? AND is_valid = TRUE`
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var sessions []*entities.Session
    for rows.Next() {
        var session entities.Session
        var refreshExpiresAt sql.NullTime
        var updatedAt sql.NullTime
        var deviceInfo sql.NullString

        err := rows.Scan(
            &session.ID,
            &session.UserID,
            &session.AccessTokenHash,
            &session.RefreshTokenHash,
            &session.ExpiresAt,
            &refreshExpiresAt,
            &session.IsValid,
            &deviceInfo,
            &session.CreatedAt,
            &updatedAt,
        )
        if err != nil {
            return nil, err
        }

        if refreshExpiresAt.Valid {
            session.RefreshExpiresAt = refreshExpiresAt.Time
        }
        if updatedAt.Valid {
            session.UpdatedAt = updatedAt.Time
        }
        if deviceInfo.Valid {
            session.DeviceInfo = deviceInfo.String
        }

        sessions = append(sessions, &session)
    }

    return sessions, nil
}