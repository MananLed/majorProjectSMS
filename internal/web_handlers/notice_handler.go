package web_handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

type NoticeHandler struct {
	Service service.NoticeServiceInterface
}

func NewNoticeHandler(s service.NoticeServiceInterface) *NoticeHandler {
	return &NoticeHandler{Service: s}
}

func (h *NoticeHandler) IssueNotice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role == model.RoleResident) {
		logger.LogToFile("unauthorized person wants to issue invoice")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Content == "" {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	}

	now := time.Now()
	if err := h.Service.IssueNotice(req.Content, now.Month(), now.Year()); err != nil {
		logger.LogToFile("Failed to issue notice: " + err.Error())
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to issue notice", 1011)
		return
	}

	logger.LogToFile("Notice issued successfully")
	response.SuccessResponse(w, nil, "Notice issued successfully", http.StatusOK)
}

func (h *NoticeHandler) GetNotices(w http.ResponseWriter, r *http.Request){
		notices, err := h.Service.GetNotices()
		if err != nil {
			logger.LogToFile("Failed to fetch notices: " + err.Error())
			response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch notices", 1011)
			return
		}
		logger.LogToFile("Notice fetched successfully")
		response.SuccessResponse(w, notices, "Notice fetched successfully", http.StatusOK)
}

func (h *NoticeHandler) GetNoticesByMonthYear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}	
	
	
	monthStr := r.URL.Query().Get("month")
	yearStr := r.URL.Query().Get("year")

	if yearStr == "" && monthStr == "" || yearStr == "" && monthStr != "" {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	} else if yearStr != "" && monthStr == "" {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			logger.LogToFile("Invalid request")
			response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
			return
		}
		invoices, err := h.Service.GetNoticesByYear(year)
		if err != nil {
			logger.LogToFile("Notices not found")
			response.ErrorResponse(w, http.StatusNotFound, "Notices not found", 404)
			return
		}
		logger.LogToFile("notices retrived successfully")
		response.SuccessResponse(w, invoices, "notices retrived successfully", http.StatusOK)
	} else {
		year, err1 := strconv.Atoi(yearStr)
		monthInt, err2 := strconv.Atoi(monthStr)
		if err1 != nil || err2 != nil || monthInt < 1 || monthInt > 12 {
			logger.LogToFile("Invalid request")
			response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
			return
		}
		month := time.Month(monthInt)
		invoice, err := h.Service.GetNoticesByMonthYear(month, year)
		if err != nil {
			logger.LogToFile("notices not found")
			response.ErrorResponse(w, http.StatusNotFound, "notices not found", 404)
			return
		}
		logger.LogToFile("notices retrived successfully")
		response.SuccessResponse(w, invoice, "notices Retrived Successfully", http.StatusOK)
	}
}

