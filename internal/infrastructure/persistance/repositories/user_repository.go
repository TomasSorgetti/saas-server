package repositories

import (
	"database/sql"
	"errors"
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

func (r *UserRepository) Save(user *entities.User) (int64, error) {
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

	return id, nil
}

func (r *UserRepository) CreateEmailVerification(userID int64, code string, expiresAt time.Time) error {
	query := `
        INSERT INTO email_verifications (user_id, verification_code, expires_at)
        VALUES (?, ?, ?)
    `
	_, err := r.db.Exec(query, userID, code, expiresAt)
	return err
}

func (r *UserRepository) UpdateEmailVerified(userID int64, verified bool) error {
	query := `
        UPDATE users SET verified = ? WHERE id = ?
    `
	_, err := r.db.Exec(query, verified, userID)
	
	return err
}

func (r *UserRepository) FindByID(id int64) (*entities.User, error) {
	query := `
        SELECT id, email, password, role, first_name, last_name, phone, address, country,
               workshop_name, is_active, deleted, last_login
        FROM users WHERE id = ?
    `
	var user entities.User
	var lastLogin sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
		&user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
		&user.Deleted, &lastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.LastLogin = lastLogin.String
	return &user, nil
}

func (r *UserRepository) FindByEmail(email string) (*entities.User, error) {
	query := `
        SELECT id, email, password, role, first_name, last_name, phone, address, country,
               workshop_name, is_active, deleted, last_login
        FROM users WHERE email = ?
    `
	var user entities.User
	var lastLogin sql.NullString

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
		&user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
		&user.Deleted, &lastLogin,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.LastLogin = lastLogin.String
	return &user, nil
}

func (r *UserRepository) FindAll() ([]*entities.User, error) {
	query := `
        SELECT id, email, password, role, first_name, last_name, phone, address, country,
               workshop_name, is_active, deleted, last_login
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
		var lastLogin sql.NullString
		if err := rows.Scan(
			&user.ID, &user.Email, &user.Password, &user.Role, &user.FirstName, &user.LastName,
			&user.Phone, &user.Address, &user.Country, &user.WorkshopName, &user.IsActive,
			&user.Deleted, &lastLogin,
		); err != nil {
			return nil, err
		}
		user.LastLogin = lastLogin.String
		users = append(users, &user)
	}
	return users, nil
}
