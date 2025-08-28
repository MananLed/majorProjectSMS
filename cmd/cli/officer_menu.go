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

// func ShowOfficerDashboard(ctx context.Context, user *model.User, uHandler *handlers.UserHandler, sHandler *handlers.ServiceRequestHandler,
// 	nHandler *handlers.NoticeHandler, fHandler *handlers.FeedbackHandler, iHandler *handlers.InvoiceHandler) {


// 	for {
// 		fmt.Println(constants.OfficerEmogiPrompt)
// 		myFigure := figure.NewColorFigure("Officer", "", "green", false)
// 		myFigure.Print()
// 		fmt.Println(constants.OfficerEmogiPrompt)
// 		color.Cyan("1. " + string(constants.ManageServiceRequestPrompt))
// 		color.Cyan("2. " + string(constants.ManageInvoicesPrompt))
// 		color.Cyan("3. " + string(constants.ManageNoticesPrompt))
// 		color.Cyan("4. " + string(constants.ManageFeedbackPrompt))
// 		color.Cyan("5. " + string(constants.ManageProfilePrompt))
// 		color.Cyan("6. " + string(constants.AddNewOfficerPrompt))
// 		color.Red("7. " + string(constants.LogoutPrompt))

// 		for {
// 			choice := utils.ReadChoice()
// 			if choice == "" {
// 				break
// 			}
// 			switch choice {
// 			case "1":
// 				menus.ShowOfficerServiceRequestMenu(ctx, sHandler)
// 			case "2":
// 				menus.ShowOfficerInvoiceMenu(iHandler)
// 			case "3":
// 				menus.ShowOfficerNoticeMenu(ctx, nHandler)
// 			case "4":
// 				menus.ShowOfficerFeedbackMenu(ctx, fHandler)
// 			case "5":
// 				isDeleted := menus.ShowOfficerResidentProfileMenu(ctx, user, uHandler)
// 				if isDeleted {
// 					return
// 				}
// 			case "6":
// 				uHandler.CreateOfficer(ctx)
// 			case "7":
// 				color.Red("Logging out...")
// 				return
// 			default:
// 				color.Red("Invalid choice, try again.")
// 			}
// 		}
// 	}
// }
