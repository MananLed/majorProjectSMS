package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"gitlab.com/david_mbuvi/go_asterisks"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService            service.UserServiceInterface
	ServiceReqeuestService *service.ServiceRequestService
}

func NewUserHandler(us service.UserServiceInterface, srs *service.ServiceRequestService) *UserHandler {
	return &UserHandler{
		UserService:            us,
		ServiceReqeuestService: srs,
	}
}

func (h *UserHandler) SignUp() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(constants.SignUpEmogiPrompt)
	myFigure := figure.NewColorFigure("Sign Up", "", "purple", false)
	myFigure.Print()
	fmt.Println(constants.SignUpEmogiPrompt)

	firstName := utils.PromptRequired(string(constants.FirstNamePrompt), reader)

	fmt.Print(color.YellowString(string(constants.MiddleNamePrompt)))
	middleName, _ := reader.ReadString('\n')
	middleName = strings.TrimSpace(middleName)

	lastName := utils.PromptRequired(string(constants.LastNamePrompt), reader)

	email := utils.PromptRequired(string(constants.EmailPrompt), reader)

	for {
		if utils.ValidateEmail(email) {
			break
		}
		email = utils.PromptRequired(string(constants.EmailPrompt), reader)
	}

	mobile := utils.PromptRequired(string(constants.MobilePrompt), reader)
	for {
		if utils.ValidateMobileNumber(mobile) {
			break
		}
		mobile = utils.PromptRequired(string(constants.MobilePrompt), reader)
	}

	id := utils.GenerateUUID().String()

	var passwordStr string
	for {
		fmt.Print(color.YellowString(string(constants.PasswordPrompt)))
		password, _ := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
		passwordStr = string(password)

		if passwordStr == "" {
			color.Red("Password is compulsory")
			continue
		}

		if !utils.ValidatePassword(passwordStr) {
			color.Red("Invalid password, enter again")
			continue
		}
		if !h.UserService.IsPasswordUnique(passwordStr) {
			color.Red("This password is already in use. Please enter a different one.")
			continue
		}
		break
	}

	var confirmPasswordStr string
	for {
		fmt.Print(color.YellowString(string(constants.ConfirmPasswordPrompt)))
		password, _ := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
		confirmPasswordStr = string(password)

		if passwordStr != confirmPasswordStr {
			color.Red("Password doesn't match")
			continue
		}
		break
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passwordStr), bcrypt.DefaultCost)

	flatNo := utils.PromptRequired(string(constants.FlatNoPrompt), reader)
	for {
		if utils.ValidateFlatNumber(flatNo) {
			break
		}
		flatNo = utils.PromptRequired(string(constants.MobilePrompt), reader)
	}

	roleStr := "FlatResident"

	user := model.User{
		FirstName:    strings.TrimSpace(firstName),
		MiddleName:   strings.TrimSpace(middleName),
		LastName:     strings.TrimSpace(lastName),
		Email:        strings.TrimSpace(email),
		MobileNumber: strings.TrimSpace(mobile),
		ID:           strings.TrimSpace(id),
		Password:     string(hashedPassword),
		Flat:         strings.TrimSpace(flatNo),
		Role:         model.ParseRole(strings.TrimSpace(roleStr)),
	}

	err := h.UserService.SignUp(user)

	if err != nil {
		color.Red("Sign up failed: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	} else {
		color.Green("User signed up successfully!!")
	}
}

func (h *UserHandler) Login() *model.User {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(constants.LoginEmogiPrompt)
	myFigure := figure.NewColorFigure("Login", "", "blue", false)
	myFigure.Print()
	fmt.Println(constants.LoginEmogiPrompt)

	var emailid string
	for {
		fmt.Print(color.YellowString(string(constants.IDPrompt)))
		input, _ := reader.ReadString('\n')
		emailid = strings.TrimSpace(input)

		if emailid == "" {
			color.Red("ID is required")
			continue
		}
		break
	}

	var passwordStr string
	for {
		fmt.Print(color.YellowString(string(constants.PasswordPrompt)))
		password, _ := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
		passwordStr = string(password)
		passwordStr = strings.TrimRight(passwordStr, "\r\n")
		if passwordStr == "" {
			color.Red("Password is required")
			continue
		}
		break
	}

	user, err := h.UserService.Login(emailid, passwordStr)
	if err != nil {
		color.Red("Login failed: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return nil
	}

	color.Green("Login successful!! Welcome, %s", user.FirstName)
	return user
}

func (h *UserHandler) UpdateProfile(user *model.User) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(color.YellowString("Update First Name: "))
	firstName, _ := reader.ReadString('\n')
	firstName = strings.TrimSpace(firstName)
	firstName = strings.TrimRight(firstName, "\r\n")
	if firstName != "" {
		user.FirstName = firstName
	}

	fmt.Print(color.YellowString("Update Middle Name: "))
	middleName, _ := reader.ReadString('\n')
	middleName = strings.TrimSpace(middleName)
	middleName = strings.TrimRight(middleName, "\r\n")
	if middleName != "" {
		user.MiddleName = middleName
	}

	fmt.Print(color.YellowString("Update Last Name: "))
	lastName, _ := reader.ReadString('\n')
	lastName = strings.TrimSpace(lastName)
	lastName = strings.TrimRight(lastName, "\r\n")
	if lastName != "" {
		user.LastName = lastName
	}

	for {
		fmt.Print(color.YellowString("Update Mobile Number: "))
		mobile, _ := reader.ReadString('\n')
		mobile = strings.TrimSpace(mobile)
		mobile = strings.TrimRight(mobile, "\r\n")
		if mobile == "" {
			break
		} else if !utils.ValidateMobileNumber(mobile) {
			color.Red("Invalid Mobile Number, enter again")
			continue
		} else {
			user.MobileNumber = mobile
			break
		}
	}

	for {
		fmt.Print(color.YellowString("Update Email: "))
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)
		email = strings.TrimRight(email, "\r\n")
		if email == "" {
			break
		} else if !utils.ValidateEmail(email) {
			color.Red("Invalid Email, enter again")
			continue
		} else {
			user.Email = email
			break
		}
	}

	if err := h.UserService.UpdateProfile(*user); err != nil {
		color.Red("Failed to update profile: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Profile updated successfully!")
}

func (h *UserHandler) ChangePassword(ctx context.Context) {

	fmt.Print(color.YellowString("Enter current password: "))
	oldPasswordBytes, _ := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
	oldPassword := string(oldPasswordBytes)
	oldPassword = strings.TrimRight(oldPassword, "\r\n")
	var newPassword, confirmNewPassword string

	for {
		fmt.Print(color.YellowString("Enter new password: "))
		passBytes, _ := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
		newPassword = string(passBytes)
		newPassword = strings.TrimRight(newPassword, "\r\n")

		if !utils.ValidatePassword(newPassword) {
			color.Red("Password must have at least 12 characters, one lowercase, one digit, and one special character.")
			continue
		}

		fmt.Print(color.YellowString("Confirm new password: "))
		confirmBytes, _ := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
		confirmNewPassword = string(confirmBytes)
		confirmNewPassword = strings.TrimRight(confirmNewPassword, "\r\n")

		if newPassword != confirmNewPassword {
			color.Red("Passwords do not match.")
			continue
		}
		break
	}

	err := h.UserService.ChangePassword(ctx, oldPassword, newPassword)
	if err != nil {
		color.Red("Failed to change password: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Password changed successfully!")
}

func (h *UserHandler) CreateOfficer(ctx context.Context) {
	currentUser, err := utils.GetUserFromContext(ctx)
	if err != nil || (currentUser.Role != model.RoleAdmin && currentUser.Role != model.RoleOfficer) {
		color.Red("Unauthorized: Only Admin or Officer can add a new officer.")
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	reader := bufio.NewReader(os.Stdin)
	email := utils.PromptRequired("Officer Email (used as ID)", reader)

	var passwordStr string
	var hashedPassword []byte

	for {
		fmt.Print(color.YellowString("Set temporary password for the officer: "))
		password, _ := go_asterisks.GetUsersPassword("", true, os.Stdin, os.Stdout)
		passwordStr = string(password)
		passwordStr = strings.TrimRight(passwordStr, "\r\n")

		if !utils.ValidatePassword(passwordStr) {
			color.Red("Invalid password format.")
			continue
		}

		hashedPassword, _ = bcrypt.GenerateFromPassword([]byte(passwordStr), bcrypt.DefaultCost)
		if !h.UserService.IsPasswordUnique(string(hashedPassword)) {
			color.Red("Password already exists. Please enter a different one.")
			continue
		}
		break
	}
	newOfficer := model.User{
		Email:        email,
		ID:           utils.GenerateUUID().String(),
		Password:     string(hashedPassword),
		Role:         model.RoleOfficer,
		FirstName:    "********",
		LastName:     "*********",
		MobileNumber: "**********",
		Flat: "xxx",
	}

	if err := h.UserService.SignUp(newOfficer); err != nil {
		color.Red("Error creating officer: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	} else {
		color.Green("Officer created successfully.")
	}
}

func (h *UserHandler) ViewProfile(user *model.User) {
	color.Cyan("\n--- User Profile ---")
	fmt.Println("First Name:", user.FirstName)
	fmt.Println("Middle Name:", user.MiddleName)
	fmt.Println("Last Name:", user.LastName)
	fmt.Println("Email/ID:", user.Email)
	fmt.Println("Mobile Number:", user.MobileNumber)
	fmt.Println("Role:", user.Role)
	if(user.Flat != "xxx") {fmt.Println("Flat:", user.Flat)}
	color.Cyan("--------------------\n")
}

func (h *UserHandler) DeleteProfile(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)
	_, err := utils.GetUserFromContext(ctx)
	if err != nil {
		color.Red("Unauthorized access.")
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	fmt.Print(color.YellowString("Are you sure you want to delete your profile? Type YES to confirm: "))
	confirmation, _ := reader.ReadString('\n')
	confirmation = strings.TrimSpace(confirmation)

	if confirmation != "YES" {
		color.Red("Profile deletion cancelled.")
		return
	}

	err = h.UserService.DeleteProfile(ctx)
	if err != nil {
		color.Red("Failed to delete profile: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	} else {
		color.Green("Profile deleted successfully!")
	}
	err = h.ServiceReqeuestService.DeleteServiceRequestByID(ctx)
	if err != nil {
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	}
}
