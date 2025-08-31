package web_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserService           service.UserServiceInterface
	ServiceRequestService *service.ServiceRequestService
}

func NewUserHandler(us service.UserServiceInterface, srs *service.ServiceRequestService) *UserHandler {
	return &UserHandler{
		UserService:           us,
		ServiceRequestService: srs,
	}
}

func (h *UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}
	var req struct {
		FirstName string `json:"firstname"`
		LastName  string `json:"lastname"`
		MiddleName string `json:"middlename"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Mobile    string `json:"mobile"`
		Flat      string `json:"flat"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	if !utils.ValidateEmail(req.Email) || !utils.ValidateMobileNumber(req.Mobile) || !utils.ValidateFlatNumber(req.Flat) || !utils.ValidatePassword(req.Password) {
		logger.LogToFile("Validation error")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	if !h.UserService.IsPasswordUnique(req.Password) {
		logger.LogToFile("Password not unique")
		response.ErrorResponse(w, http.StatusBadRequest, "Password cannot be used", 1001)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to hash password", 1010)
		return
	}

	user := model.User{
		ID:           utils.GenerateUUID().String(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		MiddleName:   req.MiddleName,
		Email:        req.Email,
		MobileNumber: req.Mobile,
		Flat:         req.Flat,
		Password:     string(hashedPassword),
		Role:         model.RoleResident,
	}

	err = h.UserService.SignUp(user)

	if err != nil {
		logger.LogToFile("Error creating user")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error creating user", 1006)
		return
	}

	logger.LogToFile("User created sucessfully")
	response.SuccessResponse(w, nil, "User created successfully", http.StatusCreated)

}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid Input")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid input", 1001)
		return
	}

	user, err := h.UserService.Login(req.Email, req.Password)
	if err != nil {
		logger.LogToFile("Invalid email or password")
		response.ErrorResponse(w, http.StatusUnauthorized, "Invalid email or password", 1005)
		return
	}

	var jwtTokenString string
	jwtTokenString, err = utils.GenerateJWT(user.ID, string(user.Role), user.Email, user.Flat)

	if err != nil {
		logger.LogToFile("Error generating the token")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error generating the token", 1006)
		return
	}

	logger.LogToFile(fmt.Sprintf("token generated for userId:%v", user.ID))
	response.SuccessResponse(w, map[string]interface{}{"token": jwtTokenString}, "Token generated Successfully", http.StatusCreated)

}

func (h *UserHandler) ViewProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	user, err := h.UserService.GetUserByID(r.Context())

	if err != nil {
		logger.LogToFile("user id not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
		return
	}

	logger.LogToFile("users retrived successfully")
	response.SuccessResponse(w, user, "Users Retrived Successfully", http.StatusOK)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}
	user, err := h.UserService.GetUserByID(r.Context())
	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	var req struct {
		FirstName    string `json:"firstName"`
		MiddleName   string `json:"middleName"`
		LastName     string `json:"lastName"`
		Email        string `json:"email"`
		MobileNumber string `json:"mobilenumber"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid JSON body")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.MiddleName != "" {
		user.MiddleName = req.MiddleName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.MobileNumber != "" {
		user.MobileNumber = req.MobileNumber
	}
	err = h.UserService.UpdateProfile(*user)

	if err != nil {
		logger.LogToFile("Failed to update user: " + err.Error())
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to update user", 1011)
		return
	}

	logger.LogToFile("User updated successfully")
	response.SuccessResponse(w, nil, "User updated successfully", http.StatusOK)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}
	
	user, err := h.UserService.GetUserByID(r.Context())
	if err != nil {
		logger.LogToFile("user not found")
        response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
        return
	}

	var req struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
        return
	}

	err = h.UserService.ChangePassword(user, req.OldPassword, req.NewPassword)
	if err != nil {
		logger.LogToFile("unauthorized user changing password")
        response.ErrorResponse(w, http.StatusUnauthorized, "User not authenticated", 1007)
        return
	}
	logger.LogToFile("User password updated successfully")
	response.SuccessResponse(w, nil, "User password updated successfully", http.StatusOK)
}

func (h *UserHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}
	ctx := r.Context()
	_, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("user not found")
        response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
        return
	}

	err = h.UserService.DeleteProfile(ctx)
	if err != nil {
		logger.LogToFile("Error deleting user")
        response.ErrorResponse(w, http.StatusInternalServerError, "Error deleting user", 1010)
        return
	}
	logger.LogToFile("Service request of deleted user deleted successfully")
    response.SuccessResponse(w, nil, "Profile deleted successfully", http.StatusOK)
}

func (h *UserHandler) CreateOfficer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role != model.RoleAdmin && currentUser.Role != model.RoleOfficer) {
		logger.LogToFile("unauthorized person wants to view profile")
        response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
        return
	}

	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
        return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to hash password", 1010)
		return
	}

	newOfficer := model.User{
		Email:    req.Email,
		ID:       utils.GenerateUUID().String(),
		Password: string(hashedPassword),
		Role:     model.RoleOfficer,
		FirstName: "********",
		LastName: "*******",
		MobileNumber: "**********",
		Flat: "xxx",
	}

	if err := h.UserService.SignUp(newOfficer); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to create officer", 1010)
		return
	}

	logger.LogToFile("Officer created successfully")
    response.SuccessResponse(w, newOfficer.ID, "Officer created successfully", http.StatusOK)
}
