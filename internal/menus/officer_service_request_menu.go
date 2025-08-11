package menus

import (
	"context"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/handlers"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/fatih/color"
)

func ShowOfficerServiceRequestMenu(ctx context.Context, sHandler *handlers.ServiceRequestHandler) {
	for {
		color.Cyan("1." + string(constants.ApproveServiceRequestPrompt))
		color.Cyan("2." + string(constants.GetPendingServiceRequestPrompt))
		color.Cyan("3." + string(constants.GetApprovedServiceRequestPrompt))
		color.Cyan("4." + "Exit")
		for {
			ch := utils.ReadChoice()
			if ch == "" {
				break
			}
			switch ch {
			case "1":
				sHandler.ApproveRequest(ctx)
			case "2":
				sHandler.ViewApprovedRequestsByServiceType(ctx)
			case "3":
				sHandler.ViewPendingRequestsByServiceType(ctx)
			case "4":
				color.Red("Exit")
				return
			default:
				color.Red("Invalid choice, try again.")
			}
		}
	}
}
