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

type InvoiceHandler struct {
	Service service.InvoiceServiceInterface
}

func NewInvoiceHandler(s service.InvoiceServiceInterface) *InvoiceHandler {
	return &InvoiceHandler{Service: s}
}

func (h *InvoiceHandler) IssueInvoice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role != model.RoleAdmin) {
		logger.LogToFile("unauthorized person wants to issue invoice")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Amount <= 0 {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	}

	now := time.Now()
	err = h.Service.GenerateInvoice(req.Amount, now.Month(), now.Year())
	if err != nil {
		logger.LogToFile("Failed to generate invoice: " + err.Error())
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to generate invoice", 1011)
		return
	}

	logger.LogToFile("Invoice issued successfully")
	response.SuccessResponse(w, nil, "Invoice issued successfully", http.StatusOK)
}

func (h *InvoiceHandler) GetInvoiceByMonthAndYear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	yearStr := r.URL.Query().Get("year")
	monthStr := r.URL.Query().Get("month")

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
		invoices, err := h.Service.GetInvoicesByYear(year)
		if err != nil {
			logger.LogToFile("invoices not found")
			response.ErrorResponse(w, http.StatusNotFound, "Invoices not found", 404)
			return
		}
		logger.LogToFile("invoices retrived successfully")
		response.SuccessResponse(w, invoices, "Invoices Retrived Successfully", http.StatusOK)
	} else {
		year, err1 := strconv.Atoi(yearStr)
		monthInt, err2 := strconv.Atoi(monthStr)
		if err1 != nil || err2 != nil || monthInt < 1 || monthInt > 12 {
			logger.LogToFile("Invalid request")
			response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
			return
		}
		month := time.Month(monthInt)
		invoice, err := h.Service.GetInvoiceByMonthAndYear(month, year)
		if err != nil {
			logger.LogToFile("invoice not found")
			response.ErrorResponse(w, http.StatusNotFound, "Invoice not found", 404)
			return
		}
		logger.LogToFile("invoice retrived successfully")
		response.SuccessResponse(w, invoice, "Invoices Retrived Successfully", http.StatusOK)
	}
}
