package main

// import (
// 	"context"
// 	"fmt"

// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/handlers"
// 	"github.com/MananLed/majorProjectSMS/internal/menus"
// 	"github.com/MananLed/majorProjectSMS/internal/model"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/common-nighthawk/go-figure"
// 	"github.com/fatih/color"
// )

// func ShowResidentDashboard(ctx context.Context, user *model.User, uHandler *handlers.UserHandler, sHandler *handlers.ServiceRequestHandler,
// 	nHandler *handlers.NoticeHandler, fHandler *handlers.FeedbackHandler, iHandler *handlers.InvoiceHandler) {

// 	for {
// 		fmt.Println(constants.ResidentEmogiPrompt)
// 		myFigure := figure.NewColorFigure("Resident", "", "green", false)
// 		myFigure.Print()
// 		fmt.Println(constants.ResidentEmogiPrompt)

// 		color.Cyan("1. " + string(constants.ManageServiceRequestPrompt))
// 		color.Cyan("2. " + string(constants.ViewNoticesPrompt))
// 		color.Cyan("3. " + string(constants.ManageFeedbackPrompt))
// 		color.Cyan("4. " + string(constants.ViewInvoicesPrompt))
// 		color.Cyan("5. " + string(constants.ManageProfilePrompt))
// 		color.Red("6. " + string(constants.LogoutPrompt))

// 		for {
// 			choice := utils.ReadChoice()
// 			if choice == "" {
// 				break
// 			}
// 			switch choice {
// 			case "1":
// 				menus.ShowResidentServiceRequestMenu(ctx, sHandler)

// 			case "2":
// 				menus.ShowResidentNoticeMenu(nHandler)

// 			case "3":
// 				menus.ShowResidentFeedbackMenu(ctx, fHandler)

// 			case "4":
// 				menus.ShowResidentInvoiceMenu(iHandler)

// 			case "5":
// 				isDeleted := menus.ShowOfficerResidentProfileMenu(ctx, user, uHandler)
// 				if isDeleted {
// 					return
// 				}

// 			case "6":
// 				color.Red("Logging out...")
// 				return

// 			default:
// 				color.Red("Invalid choice, try again.")
// 			}
// 		}
// 	}
// }
