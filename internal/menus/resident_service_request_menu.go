package menus

import (
	"context"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/handlers"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/fatih/color"
)

func ShowResidentServiceRequestMenu(ctx context.Context, sHandler *handlers.ServiceRequestHandler) {
	for {
		color.Cyan("1." + string(constants.CreateServiceRequestPrompt))
		color.Cyan("2." + string(constants.RescheduleServiceRequestPrompt))
		color.Cyan("3." + string(constants.CancelServiceRequestPrompt))
		color.Cyan("4." + string(constants.GetApprovedServiceRequestPrompt))
		color.Cyan("5." + string(constants.GetPendingServiceRequestPrompt))
		color.Cyan("6. Exit")
		for {
			ch := utils.ReadChoice()
			if ch == "" {
				break
			}
			switch ch {
			case "1":
				sHandler.BookServiceRequest(ctx)
			case "2":
				sHandler.RescheduleServiceRequest(ctx)
			case "3":
				sHandler.CancelServiceRequest(ctx)
			case "4":
				sHandler.GetApprovedServiceRequests(ctx)
			case "5":
				sHandler.GetPendingServiceRequests(ctx)
			case "6":
				color.Red("Exit")
				return
			default:
				color.Red("Invalid choice, try again.")
			}
		}
	}
}
