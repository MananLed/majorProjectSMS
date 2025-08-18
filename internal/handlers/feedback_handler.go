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
	"github.com/fatih/color"
)

type FeedbackHandler struct {
	FeedbackService service.FeedbackServiceInterface
}

func NewFeedbackHandler(service service.FeedbackServiceInterface) *FeedbackHandler {
	return &FeedbackHandler{FeedbackService: service}
}

func (h *FeedbackHandler) IssueFeedback(ctx context.Context) {

	user, err := utils.GetUserFromContext(ctx)

	reader := bufio.NewReader(os.Stdin)

	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	fmt.Print(color.YellowString("Enter Rating: "))
	var rating int32
	fmt.Scanf("%d\n", &rating)

	for {
		if rating < 1 || rating > 5 {
			color.Red("rating should be in range of 1 to 5, enter again")
		} else {
			break
		}
		fmt.Scanf("%d\n", &rating)
	}
	fmt.Print(color.YellowString("Enter Feedback: "))
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	err = h.FeedbackService.IssueFeedback(content, user.ID, rating)
	if err != nil {
		color.Red("Failed to send feedback: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	} else {
		color.Green("Feedback sent successfully.")
	}
}

func (h *FeedbackHandler) GetFeedbacks(ctx context.Context) {

	user, err := utils.GetUserFromContext(ctx)

	if user.Role == model.RoleResident {
		color.Red("not permitted to access feedbacks")
		return
	}
	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	feedbacks, err := h.FeedbackService.GetFeedbacks()

	if err != nil {
		color.Red("failed to retrieve feedbacks: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	}

	if len(feedbacks) == 0 {
		color.Yellow("No feedbacks found.")
		return
	}

	color.Cyan("===== All Feedbacks =====\n\n")

	for _, feedback := range feedbacks {
		fmt.Printf(constants.FeedbackFormatPrompt, feedback.ResidentID, feedback.Rating, feedback.Content)
	}
}

func (h *FeedbackHandler) GetFeebacksByResidentID(ctx context.Context) {
	user, err := utils.GetUserFromContext(ctx)

	if user.Role == model.RoleResident {
		color.Red("not permitted to access feedbacks")
		return
	}
	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	fmt.Print(color.YellowString("Enter Resident ID: "))
	reader := bufio.NewReader(os.Stdin)

	id, _ := reader.ReadString('\n')
	id = strings.TrimRight(id, "\r\n")
	feedbacks, err := h.FeedbackService.GetFeedbackByID(id)

	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Feedbacks of %v:-", id)
	for _, feedback := range feedbacks {
		fmt.Printf(constants.FeedbackFormatPrompt, feedback.ResidentID, feedback.Rating, feedback.Content)
	}
}

func (h *FeedbackHandler) GetFeebacksOfResident(ctx context.Context) {
	user, err := utils.GetUserFromContext(ctx)

	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	feedbacks, err := h.FeedbackService.GetFeedbackByID(user.ID)

	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Feedbacks of %v:-", user.ID)
	for _, feedback := range feedbacks {
		fmt.Printf(constants.FeedbackFormatPrompt, feedback.ResidentID, feedback.Rating, feedback.Content)
	}
}
