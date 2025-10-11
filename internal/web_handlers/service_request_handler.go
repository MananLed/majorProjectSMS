package web_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/google/uuid"
)

type ServiceRequestHandler struct {
	Service service.ServiceRequestServiceInterface
}

func NewServiceRequestHandler(s service.ServiceRequestServiceInterface) *ServiceRequestHandler {
	return &ServiceRequestHandler{Service: s}
}

func (h *ServiceRequestHandler) GetAvailableTimeSlots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	serviceType := r.URL.Query().Get("serviceType")
	if serviceType != string(model.Plumber) && serviceType != string(model.Electrician) {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	}

	slots := h.Service.GetAvailableTimeSlots(model.ServiceType(serviceType))
	logger.LogToFile("Time slot shown successfully")
	response.SuccessResponse(w, slots, "Time slot shown successfully", http.StatusOK)
}

func (h *ServiceRequestHandler) BookServiceRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	currentUser, err := utils.GetUserFromContext(r.Context())
	if err != nil || (currentUser.Role != model.RoleResident) {
		logger.LogToFile("unauthorized person wants to book request")
		response.ErrorResponse(w, http.StatusForbidden, "Unauthorized Access", 1008)
		return
	}

	var req struct {
		ServiceType string `json:"servicetype"`
		SlotID      int    `json:"slotid"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	}
	availableSlots := h.Service.GetAvailableTimeSlots(model.ServiceType(req.ServiceType))
	if req.SlotID < 1 || req.SlotID > len(availableSlots) {
		logger.LogToFile("Invalid request")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request", 1001)
		return
	}
	chosenSlot := availableSlots[req.SlotID-1]

	now := time.Now()
	formattedDate := now.Format("02-01-2006")
	fmt.Println("date: ", formattedDate)

	request := model.ServiceRequest{
		RequestID:   uuid.New(),
		ResidentID:  currentUser.ID,
		Flat:        currentUser.Flat,
		Status:      model.StatusPending,
		TimeSlot:    chosenSlot.StartTime.Format("3:04 PM") + " - " + chosenSlot.EndTime.Format("3:04 PM"),
		StartTime:   chosenSlot.StartTime,
		EndTime:     chosenSlot.EndTime,
		ServiceType: model.ServiceType(req.ServiceType),
		Date:        formattedDate,
	}

	err = h.Service.BookServiceRequest(request)

	if err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to book service request", 1010)
		return
	}

	logger.LogToFile("Service Request created successfully")
	response.SuccessResponse(w, request.RequestID, "Service Request created successfully", http.StatusOK)
}

func (h *ServiceRequestHandler) RescheduleServiceRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/service/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "reschedule" {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid path format", 1001)
		return
	}

	reqIDStr := parts[1]
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request ID", 1002)
		return
	}

	var req struct {
		Slot int `json:"slot"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}
	slotID := req.Slot

	serviceType, err := h.Service.GetServiceTypeByID(reqID)
	if err != nil {
		response.ErrorResponse(w, http.StatusNotFound, "Request not found: "+err.Error(), 1005)
		return
	}

	availableSlots := h.Service.GetAvailableTimeSlots(serviceType)
	if slotID < 1 || slotID > len(availableSlots) {
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid slot selected", 1006)
		return
	}
	chosenSlot := availableSlots[slotID-1]

	if err := h.Service.RescheduleServiceRequest(user.ID, reqID, chosenSlot, serviceType); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to reschedule: "+err.Error(), 1008)
		return
	}

	logger.LogToFile("Service Request rescheduled successfully")
	response.SuccessResponse(w, nil, "Service Request rescheduled successfully", http.StatusOK)
}

func (h *ServiceRequestHandler) CancelServiceRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/service/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 || parts[0] != "cancel" {
		logger.LogToFile("Missing request ID")
		response.ErrorResponse(w, http.StatusBadRequest, "Missing request ID", 1001)
		return
	}

	reqIDStr := parts[1]
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		logger.LogToFile("Invalid request ID")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request ID", 1002)
		return
	}

	if err := h.Service.CancelServiceRequest(user.ID, reqID); err != nil {
		logger.LogToFile("Failed to cancel the request")
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to cancel: "+err.Error(), 1008)
		return
	}

	logger.LogToFile("Service Request cancelled successfully")
	response.SuccessResponse(w, nil, "Service Request cancelled successfully", http.StatusOK)
}

func (h *ServiceRequestHandler) ApproveRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	if user.Role != "admin" && user.Role != "officer" {
		logger.LogToFile("Unauthorized Access")
		response.ErrorResponse(w, http.StatusForbidden, "Not authorized to approve requests", 1009)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/service/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "approve" {
		logger.LogToFile("Invalid Path")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid path format", 1001)
		return
	}

	reqIDStr := parts[1]
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		logger.LogToFile("Invalid Request ID")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request ID", 1002)
		return
	}

	type RequestProvider struct {
		AssignedTo string `json:"assignedto"`
	}

	var req RequestProvider
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid JSON body")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	if err := h.Service.ApproveServiceRequest(reqID, req.AssignedTo); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to approve: "+err.Error(), 1008)
		return
	}

	logger.LogToFile("Service Request approved successfully")
	response.SuccessResponse(w, nil, "Service Request approved successfully", http.StatusOK)
}

func (h *ServiceRequestHandler) CompleteRequest(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	if user.Role != "admin" && user.Role != "officer" {
		logger.LogToFile("Unauthorized Access")
		response.ErrorResponse(w, http.StatusForbidden, "Not authorized to approve requests", 1009)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/service/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[0] != "complete" {
		logger.LogToFile("Invalid Path")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid path format", 1001)
		return
	}

	reqIDStr := parts[1]
	reqID, err := uuid.Parse(reqIDStr)
	if err != nil {
		logger.LogToFile("Invalid Request ID")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request ID", 1002)
		return
	}

	if err := h.Service.CompleteServiceRequest(reqID); err != nil {
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to complete: "+err.Error(), 1008)
		return
	}

	logger.LogToFile("Service Request completed successfully")
	response.SuccessResponse(w, nil, "Service Request completed successfully", http.StatusOK)
}

func (h *ServiceRequestHandler) GetRequestsOfResident(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("user not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	status := r.URL.Query().Get("status")
	id := r.URL.Query().Get("id")
	if status == "" {
		logger.LogToFile("Missing Status")
		response.ErrorResponse(w, http.StatusBadRequest, "Missing status parameter", 1001)
		return
	}

	var requests []model.ServiceRequest

	if user.Role == model.RoleResident {
		requests = h.Service.GetServiceRequestsByStatus(user.ID, model.Status(status))
	} else {
		if id == "" {
			logger.LogToFile("Missing ID")
			response.ErrorResponse(w, http.StatusBadRequest, "Missing ID parameter", 1001)
			return
		}
		requests = h.Service.GetServiceRequestsByStatus(id, model.Status(status))
	}

	logger.LogToFile(fmt.Sprintf("Fetched %d requests with status %s", len(requests), status))
	response.SuccessResponse(w, requests, "Requests fetched successfully", http.StatusOK)
}

func (h *ServiceRequestHandler) GetRequestsByServiceTypeAndStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Method not allowed")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("User not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	if user.Role != "admin" && user.Role != "officer" {
		logger.LogToFile("Unauthorized Access")
		response.ErrorResponse(w, http.StatusForbidden, "Not authorized", 1008)
		return
	}

	status := r.URL.Query().Get("status")
	serviceType := r.URL.Query().Get("serviceType")
	if status == "" || serviceType == "" {
		logger.LogToFile("Invalid Request")
		response.ErrorResponse(w, http.StatusBadRequest, "Missing status or serviceType parameter", 1001)
		return
	}

	var requests []model.ServiceRequest

	switch {
	case serviceType == "plumber" && status == "pending":
		requests = h.Service.GetPendingRequestsByServiceType(model.Plumber)
	case serviceType == "plumber" && status == "approved":
		requests = h.Service.GetApprovedRequestsByServiceType(model.Plumber)
	case serviceType == "electrician" && status == "pending":
		requests = h.Service.GetPendingRequestsByServiceType(model.Electrician)
	case serviceType == "electrician" && status == "approved":
		requests = h.Service.GetApprovedRequestsByServiceType(model.Electrician)
	}

	logger.LogToFile("Requests fetched successfully")
	response.SuccessResponse(w, requests, "Requests fetched successfully", http.StatusOK)
}

func (h *ServiceRequestHandler) GetAllRequests(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("User not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	if user.Role != "admin" && user.Role != "officer" {
		logger.LogToFile("Unauthorized Access")
		response.ErrorResponse(w, http.StatusForbidden, "Not authorized", 1008)
		return
	}

	var pendingRequests []model.ServiceRequest
	var approvedRequests []model.ServiceRequest

	pendingRequests = h.Service.GetPendingRequestsByServiceType(model.Plumber)
	pendingRequests = append(pendingRequests, h.Service.GetPendingRequestsByServiceType(model.Electrician)...)
	approvedRequests = h.Service.GetApprovedRequestsByServiceType(model.Plumber)
	approvedRequests = append(approvedRequests, h.Service.GetApprovedRequestsByServiceType(model.Electrician)...)

	logger.LogToFile("All Requests fetched successfully.")
	
	allRequests := struct {
		Pending []model.ServiceRequest
		Approved []model.ServiceRequest
	}{
		Pending: pendingRequests,
		Approved: approvedRequests,
	}

	response.SuccessResponse(w, allRequests, "Requests fetched successfully", http.StatusOK)

}

func (h *ServiceRequestHandler) GetAllRequestsOfResident(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		logger.LogToFile("User not found")
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	if user.Role != "resident" {
		logger.LogToFile("Unauthorized Access")
		response.ErrorResponse(w, http.StatusForbidden, "Not authorized", 1008)
		return
	}
	
	var pendingRequests []model.ServiceRequest
	var approvedRequests []model.ServiceRequest
	var completedRequests []model.ServiceRequest

	pendingRequests = h.Service.GetServiceRequestsByStatus(user.ID, model.StatusPending)
	approvedRequests = h.Service.GetServiceRequestsByStatus(user.ID, model.StatusApproved)
	completedRequests = h.Service.GetServiceRequestsByStatus(user.ID, model.StatusCompleted)

	logger.LogToFile("All Requests fetched successfully.")
	
	allRequests := struct {
		Pending []model.ServiceRequest
		Approved []model.ServiceRequest
		Completed []model.ServiceRequest
	}{
		Pending: pendingRequests,
		Approved: approvedRequests,
		Completed: completedRequests,
	}

	response.SuccessResponse(w, allRequests, "Requests fetched successfully", http.StatusOK)
}