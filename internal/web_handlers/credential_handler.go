package web_handlers

import (
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type CredentialHandler struct {
	Service service.CredentialServiceInterface
}

func NewCredentialHandler(s service.CredentialServiceInterface) *CredentialHandler {
	return &CredentialHandler{Service: s}
}

func (h *CredentialHandler) DeleteOfficer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role != model.RoleAdmin) {
		logger.LogToFile("unauthorized person trying to delete officer")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		logger.LogToFile("ID missing while deleting officer")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	}

	if err := h.Service.DeleteOfficerCredentials(r.Context(), id); err != nil {
		logger.LogToFile("Error in deleting officer")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error deleting officer", 1010)
		return
	}

	logger.LogToFile("Officer deleted successfully")
	response.SuccessResponse(w, nil, "Officer deleted successfully", http.StatusOK)
}

func (h *CredentialHandler) DeleteResident(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role != model.RoleAdmin) {
		logger.LogToFile("unauthorized person trying to delete resident")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		logger.LogToFile("ID missing while deleting resident")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	}

	if err := h.Service.DeleteResidentCredentials(r.Context(), id); err != nil {
		logger.LogToFile("Error in deleting resident")
		response.ErrorResponse(w, http.StatusInternalServerError, "Error deleting resident", 1010)
		return
	}

	logger.LogToFile("Resident deleted successfully")
	response.SuccessResponse(w, nil, "Resident deleted successfully", http.StatusOK)
}
