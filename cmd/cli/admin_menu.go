package main

import (
	"context"
	"fmt"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/handlers"
	"github.com/MananLed/majorProjectSMS/internal/menus"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func ShowAdminDashboard(ctx context.Context, user *model.User, uHandler *handlers.UserHandler, socHandler *handlers.SocietyHandler, cHandler *handlers.CredentialHandler,
	nHandler *handlers.NoticeHandler, fHandler *handlers.FeedbackHandler, iHandler *handlers.InvoiceHandler, sHandler *handlers.ServiceRequestHandler) {

	for {
		fmt.Println(constants.AdminEmogiPrompt)
		myFigure := figure.NewColorFigure("Admin", "", "green", false)
		myFigure.Print()
		fmt.Println(constants.AdminEmogiPrompt)

		color.Cyan("1. " + string(constants.ManageResidentPrompt))
		color.Cyan("2. " + string(constants.ManageOfficerPrompt))
		color.Cyan("3. " + string(constants.ManageServiceRequestPrompt))
		color.Cyan("4. " + string(constants.ManageNoticesPrompt))
		color.Cyan("5. " + string(constants.ManageFeedbackPrompt))
		color.Cyan("6. " + string(constants.ManageInvoicesPrompt))
		color.Cyan("7. " + string(constants.ManageProfilePrompt))
		color.Cyan("8. " + string(constants.AddNewOfficerPrompt))
		color.Red("9." + string(constants.LogoutPrompt))
		for {
			choice := utils.ReadChoice()
			if choice == "" {
				break
			}
			switch choice {
			case "1":
				menus.ShowAdminResidentMenu(ctx, socHandler, cHandler)

			case "2":
				menus.ShowAdminOfficerMenu(ctx, socHandler, cHandler)

			case "3":
				menus.ShowAdminServiceRequestMenu(ctx, sHandler)

			case "4":
				menus.ShowAdminNoticeMenu(ctx, nHandler)

			case "5":
				menus.ShowAdminFeedbackMenu(ctx, fHandler)

			case "6":
				menus.ShowAdminInvoiceMenu(ctx, iHandler)

			case "7":
				menus.ShowAdminProfileMenu(ctx, uHandler, user)

			case "8":
				uHandler.CreateOfficer(ctx)

			case "9":
				color.Red("Logging out...")
				return

			default:
				color.Red("Invalid choice, try again.")
			}
		}
	}
}
