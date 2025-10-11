package web_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/response"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/google/uuid"
)

type FeedbackHandler struct {
	Service service.FeedbackServiceInterface
}

func NewFeedbackHandler(s service.FeedbackServiceInterface) *FeedbackHandler {
	return &FeedbackHandler{Service: s}
}

func (h *FeedbackHandler) GetFeedbacks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	residentID := r.URL.Query().Get("residentId")
	var feedbacks []model.Feedback

	if user.Role == model.RoleResident {
		feedbacks, err = h.Service.GetFeedbacks()
		if err != nil {
			logger.LogToFile("Failed to fetch feedbacks")
			response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch feedbacks: "+err.Error(), 1001)
			return
		}
	} else if user.Role == "admin" || user.Role == "officer" {
		if residentID != "" {
			feedbacks, err = h.Service.GetFeedbackByID(residentID)
		} else {
			feedbacks, err = h.Service.GetFeedbacks()
		}
		if err != nil {
			logger.LogToFile("Failed to fetch feedbacks")
			response.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch feedbacks: "+err.Error(), 1002)
			return
		}
	} else {
		logger.LogToFile("Unauthorized")
		response.ErrorResponse(w, http.StatusForbidden, "Not authorized", 1003)
		return
	}

	logger.LogToFile("Fetched feedbacks successfully")
	response.SuccessResponse(w, feedbacks, "Feedbacks fetched successfully", http.StatusOK)
}

func (h *FeedbackHandler) GiveFeedback(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.LogToFile("Invalid HTTP method")
		response.ErrorResponse(w, http.StatusMethodNotAllowed, "Method not allowed", 1000)
		return
	}

	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	if user.Role != model.RoleResident {
		response.ErrorResponse(w, http.StatusUnauthorized, "unauthorized", 1007)
		return
	}
	type FeedbackRequest struct {
		Rating  int32  `json:"rating"`
		Content string `json:"content"`
	}
	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid JSON body")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		logger.LogToFile("Invalid rating")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid rating", 1001)
		return
	}

	if len(req.Content) > 500 {
		logger.LogToFile("Content too long")
		response.ErrorResponse(w, http.StatusBadRequest, "Content length exceeded the permitted lenght", 1001)
	}

	if err := h.Service.IssueFeedback(req.Content, user.ID, user.Flat, req.Rating); err != nil {
		logger.LogToFile("Failed to issue feedback: " + err.Error())
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to issue feedback", 1011)
		return
	}
	logger.LogToFile("Feedback issued successfully")
	response.SuccessResponse(w, nil, "Feedback issued successfully", http.StatusOK)
}

func (h *FeedbackHandler) IssueFeedbackOnRequest(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	user, err := utils.GetUserFromContext(ctx)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	if user.Role != model.RoleResident {
		response.ErrorResponse(w, http.StatusUnauthorized, "unauthorized", 1007)
		return
	}

	type FeedbackRequest struct {
		Rating  int32  `json:"rating"`
		Content string `json:"content"`
		RequestID uuid.UUID `json:"requestid"`
	}

	var req FeedbackRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid JSON body")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}
	if req.Rating < 1 || req.Rating > 5 {
		logger.LogToFile("Invalid rating")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid rating", 1001)
		return
	}

	if len(req.Content) > 500 {
		logger.LogToFile("Content too long")
		response.ErrorResponse(w, http.StatusBadRequest, "Content length exceeded the permitted lenght", 1001)
	}

	if err := h.Service.IssueFeedbackOnRequest(req.Content, user.ID, user.Flat, req.Rating, req.RequestID); err != nil {
		logger.LogToFile("Failed to issue feedback: " + err.Error())
		response.ErrorResponse(w, http.StatusInternalServerError, "Failed to issue feedback", 1011)
		return
	}
	logger.LogToFile("Feedback issued successfully")
	response.SuccessResponse(w, nil, "Feedback issued successfully", http.StatusOK)
}

func(h *FeedbackHandler) IsFeedbackGiven(w http.ResponseWriter, r *http.Request){
	ctx := r.Context()
	_, err := utils.GetUserFromContext(ctx)
	if err != nil {
		response.ErrorResponse(w, http.StatusUnauthorized, "User not found", 1007)
		return
	}

	type FeedbackQuery struct{
		RequestID uuid.UUID `json:"requestid"`
	}

	var req FeedbackQuery
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.LogToFile("Invalid JSON body")
		response.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", 1001)
		return
	}

	var isGiven bool

	if isGiven, err = h.Service.IsFeedbackGiven(req.RequestID); err != nil {
		logger.LogToFile("Query Failed: " + err.Error())
		response.ErrorResponse(w, http.StatusInternalServerError, "Query Failed: ", 1011)
		return
	}

	logger.LogToFile("Query Successful");


	response.SuccessResponse(w, isGiven, "Query Successful", http.StatusOK)
}



