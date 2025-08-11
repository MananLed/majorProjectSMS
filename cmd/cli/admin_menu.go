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

func ShowAdminDashboard(ctx context.Context, user *model.User, uHandler *handlers.UserHandler, socHandler *handlers.SocietyHandler, cHandler *handlers.CredentialHandler,
	nHandler *handlers.NoticeHandler, fHandler *handlers.FeedbackHandler, iHandler *handlers.InvoiceHandler, sHandler *handlers.ServiceRequestHandler) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(constants.AdminEmogiPrompt)
		myFigure := figure.NewColorFigure("Admin", "", "green", false)
		myFigure.Print()
		fmt.Println(constants.AdminEmogiPrompt)

		color.Cyan("1." + string(constants.ManageResidentPrompt))
		color.Cyan("2. " + string(constants.ManageOfficerPrompt))
		color.Cyan("3. " + string(constants.ManageServiceRequestPrompt))
		color.Cyan("4. " + string(constants.ManageNoticesPrompt))
		color.Cyan("5. " + string(constants.ManageFeedbackPrompt))
		color.Cyan("6. " + string(constants.ManageInvoicesPrompt))
		color.Cyan("7. " + string(constants.ManageProfilePrompt))
		color.Cyan("8. " + string(constants.AddNewOfficerPrompt))
		color.Red("9." + string(constants.LogoutPrompt))
returnmain:
		fmt.Print(color.BlueString(string(constants.ChoicePrompt)))

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		if(choice == "") {continue}
		switch choice {
		case "1":
			for {
				color.Cyan("1." + string(constants.ViewResidentPrompt))
				color.Cyan("2." + string(constants.DeleteResidentPrompt))
				color.Cyan("3. Exit")
return1:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				exit := false
				if(ch == "") {continue}
				switch ch {
				case "1":
					socHandler.HandleViewResidents(ctx)
				case "2":
					cHandler.DeleteResident(ctx)
				case "3":
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
				color.Cyan("1." + string(constants.ViewOfficersPrompt))
				color.Cyan("2." + string(constants.DeleteOfficerPrompt))
				color.Cyan("3. Exit")
return2:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				exit := false
				if(ch == "") {continue}
				switch ch {
				case "1":
					socHandler.HandleViewOfficers(ctx)
				case "2":
					cHandler.DeleteOfficer(ctx)
				case "3":
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
				color.Cyan("1." + string(constants.GetPendingServiceRequestPrompt))
				color.Cyan("2." + string(constants.GetApprovedServiceRequestPrompt))
				color.Cyan("3." + string(constants.ApproveServiceRequestPrompt))
				color.Cyan("4. Exit")

return3:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				exit := false
				if(ch == "") {continue}
				switch ch {
				case "1":
					sHandler.ViewPendingRequestsByServiceType(ctx)
				case "2":
					sHandler.ViewApprovedRequestsByServiceType(ctx)
				case "3":
					sHandler.ApproveRequest(ctx)
				case "4":
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
				color.Cyan("1." + string(constants.IssueNoticePrompt))
				color.Cyan("2." + string(constants.GetNoticePrompt))
				color.Cyan("3." + string(constants.GetNoticesByMonthYear))
				color.Cyan("4." + string(constants.GetNoticesByYear))
				color.Cyan("5. Exit")
return4:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				exit := false
				if(ch == "") {continue}
				switch ch {
				case "1":
					nHandler.IssueNotice(ctx)
				case "2":
					nHandler.GetNotices()
				case "3":
					nHandler.GetNoticesByMonthYear()
				case "4":
					nHandler.GetNoticesByYear()
				case "5":
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
				color.Cyan("1." + string(constants.GetAllFeedbackPrompt))
				color.Cyan("2." + string(constants.GetFeedbacksByID))
				color.Cyan("3. Exit")
return5:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				exit := false
				if(ch == "") {continue}
				switch ch {
				case "1":
					fHandler.GetFeedbacks(ctx)
				case "2":
					fHandler.GetFeebacksByResidentID(ctx)
				case "3":
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
			for {
				color.Cyan("1." + string(constants.IssueInvoice))
				color.Cyan("2." + string(constants.SearchAInvoice))
				color.Cyan("3." + string(constants.ListInvoicesOfAYear))
				color.Cyan("4. Exit")
return6:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				exit := false
				if(ch == "") {continue}
				switch ch {
				case "1":
					iHandler.IssueInvoice(ctx)
				case "2":
					iHandler.GetInvoiceByMonthAndYear()
				case "3":
					iHandler.GetInvoicesByYear()
				case "4":
					color.Red("Exit")
					exit = true
				default:
					color.Red("Invalid choice, try again.")
				}
				if exit {
					break
				}
				goto return6
			}

		case "7":
			for {
				color.Cyan("1." + string(constants.UpdateProfilePrompt))
				color.Cyan("2." + string(constants.ChangePasswordPrompt))
				color.Cyan("3." + string(constants.ViewProfilePrompt))
				color.Cyan("4. Exit")
return7:
				fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
				ch, _ := reader.ReadString('\n')
				ch = strings.TrimSpace(ch)
				exit := false
				if(ch == "") {continue}
				switch ch {
				case "1":
					uHandler.UpdateProfile(user)
				case "2":
					uHandler.ChangePassword(ctx)
				case "3":
					uHandler.ViewProfile(user)
				case "4":
					color.Red("Exit")
					exit = true
				default:
					color.Red("Invalid choice, try again.")
				}
				if exit {
					break
				}
				goto return7 
			}

		case "8":
			uHandler.CreateOfficer(ctx)
			goto returnmain

		case "9":
			color.Red("Logging out...")
			return

		default:
			color.Red("Invalid choice, try again.")
		}
		goto returnmain
	}
}

