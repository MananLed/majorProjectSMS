package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	mu sync.Mutex
	db *sql.DB
}

type UserRepositoryInterface interface {
	AddUser(user model.User) error
	GetUserByIDAndPassword(id string, password string) (*model.User, error)
	UpdateUser(user model.User) error
	ChangePassword(id string, newHashedPassword string) error
	IsPasswordUnique(hashedPassword string) bool
	DeleteUserByID(id string) error
	GetUserByID(id string) (*model.User, error)
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) AddUser(newUser model.User) error {
	
	query := `
		INSERT INTO users (id, first_name, middle_name, last_name, mobile_number, email, password, role, flat_no)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query, newUser.ID, newUser.FirstName, newUser.MiddleName, newUser.LastName, newUser.MobileNumber, newUser.Email, newUser.Password, newUser.Role, newUser.Flat)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByIDAndPassword(email string, password string) (*model.User, error) {
	query := `
		SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role, flat_no
		FROM users
		WHERE email = $1
	`


	rows, err := r.db.Query(query, email)

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user model.User
		err := rows.Scan(
			&user.ID, &user.FirstName, &user.MiddleName, &user.LastName,
			&user.MobileNumber, &user.Email, &user.Password, &user.Role, &user.Flat,
		)
		if err != nil {
			logger.LogToFile(fmt.Sprintf("Row scan error: %v", err))
			continue
		}

		if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil {
			return &user, nil
		}
	}

	return nil, errors.New("user not found or password incorrect")
}

func (r *UserRepository) UpdateUser(updatedUser model.User) error {
	
	query := `
		UPDATE users
		SET first_name = $1, middle_name = $2, last_name = $3, mobile_number = $4, email = $5, password = $6, role = $7
		WHERE id = $8
	`

	result, err := r.db.Exec(query, updatedUser.FirstName, updatedUser.MiddleName, updatedUser.LastName, updatedUser.MobileNumber, updatedUser.Email, updatedUser.Password, updatedUser.Role, updatedUser.ID)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("UpdateUser error: %v", err))
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) ChangePassword(id string, newHashedPassword string) error {
	
	query := `
		UPDATE users
		SET password = $1
		WHERE id = $2
	`

	result, err := r.db.Exec(query, newHashedPassword, id)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *UserRepository) IsPasswordUnique(password string) bool {

	rows, err := r.db.Query(`
		SELECT password FROM users
	`)



	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return false
	}

	defer rows.Close()

	for rows.Next() {
		var hashed string
		if err := rows.Scan(&hashed); err != nil {
			return false
		}
		if bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil {
			return false
		}
	}
	return true
}

func (r *UserRepository) DeleteUserByID(id string) error {
	
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := r.db.Exec(query, id)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *UserRepository) GetUserByID(id string) (*model.User, error) {
	var user model.User 

	query := `
		SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role, flat_no
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email, &user.Password, &user.Role, &user.Flat)


	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		logger.LogToFile(fmt.Sprintf("error fetching user: %v", err))
		return nil, err
	}

	return &user, nil
}