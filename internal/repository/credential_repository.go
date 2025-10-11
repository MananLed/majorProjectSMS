package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type CredentialRepository struct {
	mu sync.Mutex
	db *sql.DB
}

type CredentialRepositoryInterface interface {
	DeleteUserByIDAndRole(id string, role model.UserRole) error
}

func NewCredentialRepository(db *sql.DB) *CredentialRepository {
	return &CredentialRepository{db: db}
}

func (r *CredentialRepository) DeleteUserByIDAndRole(id string, role model.UserRole) error {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = $2)`


	err := r.db.QueryRow(query, id, string(role)).Scan(&exists)


	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return fmt.Errorf("failed to check user existence: %v", err)
	}

	if !exists {
		return errors.New("user not found")
	}


	query = `DELETE FROM users WHERE id = $1 AND role = $2`


	_, err = r.db.Exec(query, id, string(role))

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}