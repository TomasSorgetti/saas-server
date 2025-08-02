package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"luthierSaas/internal/domain/entities"
	"time"

	"github.com/go-sql-driver/mysql"
)

type UserRepository struct {
	db *sql.DB
}

var ErrEmailAlreadyExists = errors.New("email already exists")

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(user *entities.User) (int, error) {
	query := `
        INSERT INTO users (
            email, password, role, first_name, last_name, phone, address, country,
            workshop_name, is_active, deleted, last_login
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	res, err := r.db.Exec(query,
		user.Email, user.Password, user.Role, user.FirstName, user.LastName,
		user.Phone, user.Address, user.Country, user.WorkshopName, user.IsActive,
		user.Deleted, sql.NullString{String: user.LastLogin, Valid: user.LastLogin != ""},
	)

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return 0, ErrEmailAlreadyExists
		}
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *UserRepository) CreateEmailVerification(userID int, code string, expiresAt time.Time) error {
	query := `
        INSERT INTO email_verifications (user_id, verification_code, expires_at)
        VALUES (?, ?, ?)
    `
	_, err := r.db.Exec(query, userID, code, expiresAt)
	return err
}

func (r *UserRepository) UpdateEmailVerified(userID int, verified bool) error {
	query := `
        UPDATE users SET verified = ? WHERE id = ?
    `
	_, err := r.db.Exec(query, verified, userID)
	
	return err
}

func (r *UserRepository) FindByID(id int) (*entities.User, error) {
    query := `
        SELECT 
            u.id, u.email, u.password, u.role, u.first_name, u.last_name, u.phone, 
            u.address, u.country, u.workshop_name, u.is_active, u.deleted, u.last_login, u.verified, u.login_method,
            s.id, s.user_id, s.plan_id, sp.name, s.status, s.started_at, s.expires_at
        FROM users u
        LEFT JOIN subscriptions s ON u.id = s.user_id AND s.status = 'active'
        LEFT JOIN subscription_plans sp ON s.plan_id = sp.id
        WHERE u.id = ?
    `
    var user entities.User
    var lastLogin sql.NullString
    var subID, subUserID, subPlanID sql.NullInt64
    var subPlanName, subStatus sql.NullString
    var subStartedAt, subExpiresAt sql.NullTime

    err := r.db.QueryRow(query, id).Scan(
        &user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
        &user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
        &user.Deleted, &lastLogin, &user.Verified, &user.LoginMethod,
        &subID, &subUserID, &subPlanID, &subPlanName, &subStatus, &subStartedAt, &subExpiresAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to query user by ID %d: %w", id, err)
    }

    user.LastLogin = lastLogin.String
    if subID.Valid {
        user.Subscription = &entities.Subscription{
            ID:        int(subID.Int64),
            UserID:    int(subUserID.Int64),
            PlanID:    int(subPlanID.Int64),
            PlanName:  subPlanName.String,
            Status:    subStatus.String,
            StartedAt: subStartedAt.Time,
            ExpiresAt: subExpiresAt.Time,
        }
    }

    return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
    query := `
        SELECT 
            u.id, u.email, u.password, u.role, u.first_name, u.last_name, u.phone, 
            u.address, u.country, u.workshop_name, u.is_active, u.deleted, u.last_login, u.verified, u.login_method,
            s.id, s.user_id, s.plan_id, sp.name, s.status, s.started_at, s.expires_at
        FROM users u
        LEFT JOIN subscriptions s ON u.id = s.user_id AND s.status = 'active'
        LEFT JOIN subscription_plans sp ON s.plan_id = sp.id
        WHERE u.email = ?
    `
    var user entities.User
    var lastLogin sql.NullString
    var subID, subUserID, subPlanID sql.NullInt64
    var subPlanName, subStatus sql.NullString
    var subStartedAt, subExpiresAt sql.NullTime

    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
        &user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
        &user.Deleted, &lastLogin, &user.Verified, &user.LoginMethod,
        &subID, &subUserID, &subPlanID, &subPlanName, &subStatus, &subStartedAt, &subExpiresAt,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, fmt.Errorf("failed to query user by email %q: %w", email, err)
    }

    user.LastLogin = lastLogin.String
    if subID.Valid {
        user.Subscription = &entities.Subscription{
            ID:        int(subID.Int64),
            UserID:    int(subUserID.Int64),
            PlanID:    int(subPlanID.Int64),
            PlanName:  subPlanName.String,
            Status:    subStatus.String,
            StartedAt: subStartedAt.Time,
            ExpiresAt: subExpiresAt.Time,
        }
    }

    return &user, nil
}

func (r *UserRepository) FindAll() ([]*entities.User, error) {
    query := `
        SELECT 
            u.id, u.email, u.password, u.role, u.first_name, u.last_name, u.phone, 
            u.address, u.country, u.workshop_name, u.is_active, u.deleted, u.last_login, u.verified, u.login_method,
            s.id, s.user_id, s.plan_id, sp.name, s.status, s.started_at, s.expires_at
        FROM users u
        LEFT JOIN subscriptions s ON u.id = s.user_id AND s.status = 'active'
        LEFT JOIN subscription_plans sp ON s.plan_id = sp.id
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("failed to query all users: %w", err)
    }
    defer rows.Close()

    var users []*entities.User
    for rows.Next() {
        var user entities.User
        var lastLogin sql.NullString
        var subID, subUserID, subPlanID sql.NullInt64
        var subPlanName, subStatus sql.NullString
        var subStartedAt, subExpiresAt sql.NullTime

        if err := rows.Scan(
            &user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
            &user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
            &user.Deleted, &lastLogin, &user.Verified, &user.LoginMethod,
            &subID, &subUserID, &subPlanID, &subPlanName, &subStatus, &subStartedAt, &subExpiresAt,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }

        user.LastLogin = lastLogin.String
        if subID.Valid {
            user.Subscription = &entities.Subscription{
                ID:        int(subID.Int64),
                UserID:    int(subUserID.Int64),
                PlanID:    int(subPlanID.Int64),
                PlanName:  subPlanName.String,
                Status:    subStatus.String,
                StartedAt: subStartedAt.Time,
                ExpiresAt: subExpiresAt.Time,
            }
        }

        users = append(users, &user)
    }

    return users, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
    query := `SELECT EXISTS (SELECT 1 FROM users WHERE email = ?)`
    var exists bool
    err := r.db.QueryRow(query, email).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to query email exists: %w", err)
    }
    return exists, nil
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID int, lastLogin time.Time) error {
    query := `UPDATE users SET last_login = ? WHERE id = ?`
    _, err := r.db.ExecContext(ctx, query, lastLogin, userID)
    return err
}