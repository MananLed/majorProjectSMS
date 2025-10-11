package web_handlers

import (
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type SocietyHandler struct {
	SocietyService service.SocietyServiceInterface
}

func NewSocietyHandler(s service.SocietyServiceInterface) *SocietyHandler {
	return &SocietyHandler{SocietyService: s}
}

func (h *SocietyHandler) ViewResidents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role != model.RoleAdmin) {
		logger.LogToFile("unauthorized person wants to view list of residents")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	residents, err := h.SocietyService.GetAllResidents(r.Context())

	if err != nil {
		logger.LogToFile("Failed to fetch residents")
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch residents", 1010)
		return
	}

	logger.LogToFile("View residents successfully executed")
	response.SuccessResponse(w, residents, "", http.StatusOK)
}

func (h *SocietyHandler) ViewOfficers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role != model.RoleAdmin) {
		logger.LogToFile("unauthorized person wants to view list of officers")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	officers, err := h.SocietyService.GetAllOfficers(r.Context())

	if err != nil {
		logger.LogToFile("Failed to fetch officers")
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch officers", 1010)
		return
	}

	logger.LogToFile("View officers successfully executed")
	response.SuccessResponse(w, officers, "", http.StatusOK)
}

func (h *SocietyHandler) GetResidentCount(w http.ResponseWriter, r *http.Request){
	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role == model.RoleResident) {
		logger.LogToFile("unauthorized access to resident count")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	residents, err := h.SocietyService.GetAllResidents(r.Context())

	if err != nil {
		logger.LogToFile("Failed to fetch residents count")
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch residents count", 1010)
		return
	}

	var countOfResident int = len(residents)

	logger.LogToFile("Count of residents successfully executed")
	response.SuccessResponse(w, countOfResident, "", http.StatusOK)
}

func (h *SocietyHandler) GetOfficerCount(w http.ResponseWriter, r *http.Request){
	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role == model.RoleResident) {
		logger.LogToFile("unauthorized access to officer count")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	officers, err := h.SocietyService.GetAllOfficers(r.Context())

	if err != nil {
		logger.LogToFile("Failed to fetch officers count")
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch officers count", 1010)
		return
	}

	var countOfOfficers int = len(officers)

	logger.LogToFile("Count of officers successfully executed")
	response.SuccessResponse(w, countOfOfficers, "", http.StatusOK)
}