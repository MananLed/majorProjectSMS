package model

type UserRole string

const(
	RoleAdmin UserRole = "admin"
	RoleOfficer UserRole = "officer"
	RoleResident UserRole = "resident"
)

type User struct{
	FirstName string
	MiddleName string 
	LastName string 
	MobileNumber string
	Email string 
	ID string  
	Password string
	Role UserRole
}


func ParseRole(role string) UserRole{
	switch role{
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