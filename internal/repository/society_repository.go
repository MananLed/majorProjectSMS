package repository

import (
	"database/sql"
	"sync"

	"fmt"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/fatih/color"
)

type SocietyRepository struct {
	mu sync.Mutex
	db *sql.DB
}

type SocietyRepositoryInterface interface {
	GetAllResidents() ([]model.User, error)
	GetAllOfficers() ([]model.User, error)
}

func NewSocietyRepository(db *sql.DB) *SocietyRepository {
	return &SocietyRepository{db: db}
}

func (s *SocietyRepository) GetAllResidents() ([]model.User, error) {

	query := `
		SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		FROM users WHERE role = $1
	`

	s.mu.Lock()
	rows, err := s.db.Query(query, string(model.RoleResident))
	s.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, fmt.Errorf("failed to fetch residents: %w", err)
	}
	defer rows.Close()

	var residents []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email, &user.Password, &user.Role); err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, fmt.Errorf("failed to scan resident: %w", err)
		}
		residents = append(residents, user)
	}

	if len(residents) > 0 {
		fmt.Print(color.YellowString("Total Residents: "), len(residents))
	} else {
		fmt.Print("There are no residents currently.")
	}
	fmt.Println()

	return residents, nil
}

func (s *SocietyRepository) GetAllOfficers() ([]model.User, error) {

	query := `
		SELECT id, first_name, middle_name, last_name, mobile_number, email, password, role 
		FROM users WHERE role = $1
	`

	s.mu.Lock()
	rows, err := s.db.Query(query, string(model.RoleOfficer))
	s.mu.Unlock()

	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil, fmt.Errorf("failed to fetch officers: %v", err)
	}
	defer rows.Close()

	var officers []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.MobileNumber, &user.Email, &user.Password, &user.Role); err != nil {
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			return nil, fmt.Errorf("failed to scan resident: %w", err)
		}
		officers = append(officers, user)
	}

	if len(officers) > 0 {
		fmt.Print(color.YellowString("Total Officers: "), len(officers))
	} else {
		fmt.Print("There are no officers currently.")
	}
	fmt.Println()

	return officers, nil
}