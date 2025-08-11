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

type ServiceRequestHandler struct {
	ServiceRequestService *service.ServiceRequestService
}

func NewServiceRequestHandler(service *service.ServiceRequestService) *ServiceRequestHandler {
	return &ServiceRequestHandler{ServiceRequestService: service}
}

func (h *ServiceRequestHandler) BookServiceRequest(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(color.YellowString("Enter service type (electrician/plumber): "))

	serviceType, _ := reader.ReadString('\n')
	serviceType = strings.ToLower(strings.TrimSpace(serviceType))

	for {
		if serviceType != string(model.Electrician) && serviceType != string(model.Plumber) {
			color.Red("Invalid service type, enter correct service type")
		} else {
			break
		}
		serviceType, _ = reader.ReadString('\n')
		serviceType = strings.ToLower(strings.TrimSpace(serviceType))
	}

	availableSlots := h.ServiceRequestService.GetAvailableTimeSlots(model.ServiceType(serviceType))

	if len(availableSlots) == 0 {
		color.Red("No available time slots.")
		return
	}

	chosenSlot, err := utils.SlotBooking(availableSlots)
	for {
		if err != nil {
			color.Red(fmt.Sprint(err))
			logger.LogToFile(fmt.Sprintf("error: %v", err))
			chosenSlot, err = utils.SlotBooking(availableSlots)
		} else {
			break
		}
	}

	requestID := utils.GenerateUUID()

	user, _ := utils.GetUserFromContext(ctx)

	request := model.ServiceRequest{
		RequestID:   requestID,
		ResidentID:  user.ID,
		Status:      model.StatusPending,
		TimeSlot:    fmt.Sprintf("%s to %s", chosenSlot.StartTime.Format("3:04 PM"), chosenSlot.EndTime.Format("3:04 PM")),
		StartTime:   chosenSlot.StartTime,
		EndTime:     chosenSlot.EndTime,
		ServiceType: model.ServiceType(serviceType),
	}

	if err := h.ServiceRequestService.BookServiceRequest(request); err != nil {
		color.Red("Error booking service request: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Service request booked successfully!")
}

func (h *ServiceRequestHandler) RescheduleServiceRequest(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(color.YellowString("Enter Request ID to reschedule: "))
	reqIDInput, _ := reader.ReadString('\n')
	reqID := strings.TrimSpace(reqIDInput)

	user, _ := utils.GetUserFromContext(ctx)

	serviceType, err := h.ServiceRequestService.GetServiceTypeByID(reqID)

	if err != nil {
		color.Red("Error: User with such ID does not exist")
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	}

	availableSlots := h.ServiceRequestService.GetAvailableTimeSlots(serviceType)

	if len(availableSlots) == 0 {
		color.Red("No available time slots.")
		return
	}

	chosenSlot, err := utils.SlotBooking(availableSlots)
	for {
		if err != nil {
			color.Red(fmt.Sprint(err))
			chosenSlot, err = utils.SlotBooking(availableSlots)
			logger.LogToFile(fmt.Sprintf("error: %v", err))
		} else {
			break
		}
	}

	err = h.ServiceRequestService.RescheduleServiceRequest(user.ID, reqID, chosenSlot, serviceType)
	if err != nil {
		color.Red("Error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Service request rescheduled successfully!")
}

func (h *ServiceRequestHandler) CancelServiceRequest(ctx context.Context) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(color.YellowString("Enter Request ID to cancel: "))

	reqIDInput, _ := reader.ReadString('\n')
	reqID := strings.TrimSpace(reqIDInput)

	user, _ := utils.GetUserFromContext(ctx)

	err := h.ServiceRequestService.CancelServiceRequest(user.ID, reqID)
	if err != nil {
		color.Red("Failed to cancel request: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Request cancelled successfully.")
}

func (h *ServiceRequestHandler) GetPendingServiceRequests(ctx context.Context) {
	h.printRequestsByStatus(ctx, model.StatusPending)
}

func (h *ServiceRequestHandler) GetApprovedServiceRequests(ctx context.Context) {
	h.printRequestsByStatus(ctx, model.StatusApproved)
}

func (h *ServiceRequestHandler) printRequestsByStatus(ctx context.Context, status model.Status) {
	user, _ := utils.GetUserFromContext(ctx)

	requests := h.ServiceRequestService.GetServiceRequestsByStatus(user.ID, status)

	if len(requests) == 0 {
		color.Yellow("No service requests with status: %s", status)
		return
	}

	for _, r := range requests {
		color.White(string(constants.SerivceFormatPrompt), r.RequestID, r.ServiceType, r.ResidentID, r.TimeSlot, r.Status)
	}
}

func (h *ServiceRequestHandler) ViewPendingRequestsByServiceType(ctx context.Context) {
	user, _ := utils.GetUserFromContext(ctx)
	if user.Role == model.RoleResident {
		color.Red("unauthorized: only officers and admins can access this")
		return
	}
	fmt.Print(color.YellowString("Enter the service for which you want to see pending requests (electrician/plumber): "))
	reader := bufio.NewReader(os.Stdin)
	serviceType, _ := reader.ReadString('\n')
	serviceType = strings.ToLower(strings.TrimSpace(serviceType))

	for {
		if serviceType != "electrician" && serviceType != "plumber" {
			color.Red("invalid service type, enter valid service type")
		} else {
			break
		}
		serviceType, _ = reader.ReadString('\n')
		serviceType = strings.ToLower(strings.TrimSpace(serviceType))
	}

	var requests []model.ServiceRequest

	switch serviceType {
	case "electrician":
		requests = h.ServiceRequestService.GetPendingRequestsByServiceType(model.Electrician)
	case "plumber":
		requests = h.ServiceRequestService.GetPendingRequestsByServiceType(model.Plumber)
	}

	if len(requests) == 0 {
		color.Yellow("No pending requests.")
		return
	}

	for _, r := range requests {
		color.White(string(constants.SerivceFormatPrompt), r.RequestID, r.ServiceType, r.ResidentID, r.TimeSlot, r.Status)
	}
}

func (h *ServiceRequestHandler) ViewApprovedRequestsByServiceType(ctx context.Context) {
	user, _ := utils.GetUserFromContext(ctx)
	if user.Role == model.RoleResident {
		color.Red("unauthorized: only officers and admins can access this")
		return
	}
	fmt.Print(color.YellowString("Enter the service for which you want to see pending requests (electrician/plumber): "))
	reader := bufio.NewReader(os.Stdin)
	serviceType, _ := reader.ReadString('\n')
	serviceType = strings.ToLower(strings.TrimSpace(serviceType))

	for {
		if serviceType != "electrician" && serviceType != "plumber" {
			color.Red("invalid service type, enter valid service type")
		} else {
			break
		}
		serviceType, _ = reader.ReadString('\n')
		serviceType = strings.ToLower(strings.TrimSpace(serviceType))
	}

	var requests []model.ServiceRequest

	switch serviceType {
	case "electrician":
		requests = h.ServiceRequestService.GetApprovedRequestsByServiceType(model.Electrician)
	case "plumber":
		requests = h.ServiceRequestService.GetApprovedRequestsByServiceType(model.Plumber)
	}

	if len(requests) == 0 {
		color.Yellow("No approved requests.")
		return
	}

	for _, r := range requests {
		color.White(string(constants.SerivceFormatPrompt), r.RequestID, r.ServiceType, r.ResidentID, r.TimeSlot, r.Status)
	}
}

func (h *ServiceRequestHandler) ApproveRequest(ctx context.Context) {
	user, _ := utils.GetUserFromContext(ctx)
	if user.Role == model.RoleResident {
		color.Red("unauthorized: Only officers and admin can approve requests")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("Enter Request ID to approve: "))
	reqIDInput, _ := reader.ReadString('\n')
	reqID := strings.TrimSpace(reqIDInput)

	err := h.ServiceRequestService.ApproveServiceRequest(reqID)
	if err != nil {
		color.Red("Failed to approve request: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Request approved successfully.")
}
