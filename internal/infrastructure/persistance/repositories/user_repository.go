package repositories

import (
	"database/sql"
	"errors"
	"luthierSaas/internal/domain/entities"

	"github.com/go-sql-driver/mysql"
)

type MySQLUserRepository struct {
    db *sql.DB
}

var ErrEmailAlreadyExists = errors.New("email allready exists")

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
    return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Save(user *entities.User) error {
    query := `
        INSERT INTO users (
            email, password, role, first_name, last_name, phone, address, country,
            workshop_name, is_active, deleted, last_login, subscription_plan,
            subscription_status, reset_password_token, reset_password_expires
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
    _, err := r.db.Exec(query,
        user.Email, user.Password, user.Role, user.FirstName, user.LastName,
        user.Phone, user.Address, user.Country, user.WorkshopName, user.IsActive,
        user.Deleted, sql.NullString{String: user.LastLogin, Valid: user.LastLogin != ""},
        user.SubscriptionPlan, user.SubscriptionStatus,
        sql.NullString{String: user.ResetPasswordToken, Valid: user.ResetPasswordToken != ""},
        sql.NullString{String: user.ResetPasswordExpires, Valid: user.ResetPasswordExpires != ""},
    )

    if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			return ErrEmailAlreadyExists
		}
		return err
	}

	return nil
}

func (r *MySQLUserRepository) FindByID(id int64) (*entities.User, error) {
    query := `
        SELECT id, email, password, role, first_name, last_name, phone, address, country,
               workshop_name, is_active, deleted, last_login, subscription_plan,
               subscription_status, reset_password_token, reset_password_expires
        FROM users WHERE id = ?
    `
    var user entities.User
    var lastLogin, resetPasswordToken, resetPasswordExpires sql.NullString
    err := r.db.QueryRow(query, id).Scan(
        &user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
        &user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
        &user.Deleted, &lastLogin, &user.SubscriptionPlan, &user.SubscriptionStatus,
        &resetPasswordToken, &resetPasswordExpires,
    )
    if err == sql.ErrNoRows {
        return nil, nil 
    }
    if err != nil {
        return nil, err
    }
    user.LastLogin = lastLogin.String
    user.ResetPasswordToken = resetPasswordToken.String
    user.ResetPasswordExpires = resetPasswordExpires.String
    return &user, nil
}

func (r *MySQLUserRepository) FindByEmail(email string) (*entities.User, error) {
    query := `
        SELECT id, email, password, role, first_name, last_name, phone, address, country,
               workshop_name, is_active, deleted, last_login, subscription_plan,
               subscription_status, reset_password_token, reset_password_expires
        FROM users WHERE email = ?
    `
    var user entities.User
    var lastLogin, resetPasswordToken, resetPasswordExpires sql.NullString
    err := r.db.QueryRow(query, email).Scan(
        &user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
        &user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
        &user.Deleted, &lastLogin, &user.SubscriptionPlan, &user.SubscriptionStatus,
        &resetPasswordToken, &resetPasswordExpires,
    )
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }
    user.LastLogin = lastLogin.String
    user.ResetPasswordToken = resetPasswordToken.String
    user.ResetPasswordExpires = resetPasswordExpires.String
    return &user, nil
}

func (r *MySQLUserRepository) FindAll() ([]*entities.User, error) {
    query := `
        SELECT id, email, password, role, first_name, last_name, phone, address, country,
               workshop_name, is_active, deleted, last_login, subscription_plan,
               subscription_status, reset_password_token, reset_password_expires
        FROM users
    `
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []*entities.User
    for rows.Next() {
        var user entities.User
        var lastLogin, resetPasswordToken, resetPasswordExpires sql.NullString
        if err := rows.Scan(
            &user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
            &user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
            &user.Deleted, &lastLogin, &user.SubscriptionPlan, &user.SubscriptionStatus,
            &resetPasswordToken, &resetPasswordExpires,
        ); err != nil {
            return nil, err
        }
        user.LastLogin = lastLogin.String
        user.ResetPasswordToken = resetPasswordToken.String
        user.ResetPasswordExpires = resetPasswordExpires.String
        users = append(users, &user)
    }
    return users, nil
}