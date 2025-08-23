package model


type UserRole string

const (
	RoleAdmin    UserRole = "admin"
	RoleOfficer  UserRole = "officer"
	RoleResident UserRole = "resident"
)

type User struct {
	FirstName    string    `json:"first_name"`
	MiddleName   string    `json:"middle_name"`
	LastName     string    `json:"last_name"`
	MobileNumber string    `json:"mobile_number"`
	Email        string    `json:"email"`
	ID           string    `json:"id"`
	Flat         string    `json:"flat"`
	Password     string    `json:"password"`
	Role         UserRole  `json:"role"`
}

func ParseRole(role string) UserRole {
	switch role {
	case "Admin":
		return RoleAdmin
	case "MaintenanceOfficer":
		return RoleOfficer
	case "FlatResident":
		return RoleResident
	default:
		return RoleResident
	}
}
