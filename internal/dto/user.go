package dto

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpRequestDTO struct {
	FirstName  string `json:"firstname"`
	LastName   string `json:"lastname"`
	MiddleName string `json:"middlename"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Mobile     string `json:"mobile"`
	Flat       string `json:"flat"`
}

type OfficerDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ChangePassword struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

type UpdateProfile struct {
	FirstName    string `json:"firstName"`
	MiddleName   string `json:"middleName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobilenumber"`
}
