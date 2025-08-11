package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/handlers"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

func ShowResidentDashboard(ctx context.Context, user *model.User, uHandler *handlers.UserHandler, sHandler *handlers.ServiceRequestHandler,
	nHandler *handlers.NoticeHandler, fHandler *handlers.FeedbackHandler, iHandler *handlers.InvoiceHandler) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(constants.ResidentEmogiPrompt)
		myFigure := figure.NewColorFigure("Resident", "", "green", false)
		myFigure.Print()
		fmt.Println(constants.ResidentEmogiPrompt)

		color.Cyan("1. " + string(constants.ManageServiceRequestPrompt))
		color.Cyan("2. " + string(constants.ViewNoticesPrompt))
		color.Cyan("3. " + string(constants.ManageFeedbackPrompt))
		color.Cyan("4. " + string(constants.ViewInvoicesPrompt))
		color.Cyan("5. " + string(constants.ManageProfilePrompt))
		color.Red("6. " + string(constants.LogoutPrompt))

returnmain:
		fmt.Print(color.BlueString(string(constants.ChoicePrompt)))

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		if(choice == "") {continue}
		switch choice {
		case "1":
			for {
				color.Cyan("1." + string(constants.CreateServiceRequestPrompt))
				color.Cyan("2." + string(constants.RescheduleServiceRequestPrompt))
				color.Cyan("3." + string(constants.CancelServiceRequestPrompt))
				color.Cyan("4." + string(constants.GetApprovedServiceRequestPrompt))
				color.Cyan("5." + string(constants.GetPendingServiceRequestPrompt))
				color.Cyan("6. Exit")

return1:				
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				if(ch == "") {continue}
				exit := false
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
					exit = true
				default:
					color.Red("Invalid choice, try again.")
				}
				if exit {
					break
				}
				goto return1
			}

		case "2":
			for {
				color.Cyan("1." + string(constants.GetNoticePrompt))
				color.Cyan("2." + string(constants.GetNoticesByMonthYear))
				color.Cyan("3." + string(constants.GetNoticesByYear))
				color.Cyan("4. Exit")

return2:				
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				if(ch == "") {continue}
				exit := false

				switch ch {
				case "1":
					nHandler.GetNotices()
				case "2":
					nHandler.GetNoticesByMonthYear()
				case "3":
					nHandler.GetNoticesByYear()
				case "4":
					color.Red("Exit")
					exit = true
				default:
					color.Red("Invalid choice, try again.")
				}
				if exit {
					break
				}
				goto return2
			}

		case "3":
			for {
				color.Cyan("1." + string(constants.GetFeedbackForResident))
				color.Cyan("2." + string(constants.GiveFeedbackPrompt))
				color.Cyan("3. Exit")

return3:
				
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				if(ch == "") {continue}
				exit := false

				switch ch {
				case "1":
					fHandler.GetFeebacksOfResident(ctx)
				case "2":
					fHandler.IssueFeedback(ctx)
				case "3":
					color.Red("Exit")
					exit = true
				default:
					color.Red("Invalid choice, try again.")
				}
				if exit {
					break
				}
				goto return3
			}

		case "4":
			for {
				color.Cyan("1." + string(constants.SearchAInvoice))
				color.Cyan("2." + string(constants.ListInvoicesOfAYear))
				color.Cyan("3. Exit")
return4:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				if(ch == "") {continue}
				exit := false

				switch ch {
				case "1":
					iHandler.GetInvoiceByMonthAndYear()
				case "2":
					iHandler.GetInvoicesByYear()
				case "3":
					color.Red("Exit")
					exit = true
				default:
					color.Red("Invalid choice, try again.")
				}
				if exit {
					break
				}
				goto return4
			}

		case "5": 
			for {
				color.Cyan("1." + string(constants.UpdateProfilePrompt))
				color.Cyan("2." + string(constants.ChangePasswordPrompt))
				color.Cyan("3." + string(constants.ViewProfilePrompt))
				color.Red("4." + string(constants.DeleteProfilePrompt))
				color.Cyan("5. Exit")
return5:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				if(ch == "") {continue}
				exit := false

				switch ch {
				case "1":
					uHandler.UpdateProfile(user)
				case "2":
					uHandler.ChangePassword(ctx)
				case "3":
					uHandler.ViewProfile(user)
				case "4":
					uHandler.DeleteProfile(ctx)
					return
				case "5":
					color.Red("Exit")
					exit = true
				default:
					color.Red("Invalid choice, try again.")
				}
				if exit {
					break
				}
				goto return5
			}

		case "6": 
			color.Red("Logging out...")
			return

		default:
			color.Red("Invalid choice, try again.")
		}
		goto returnmain
	}
}
