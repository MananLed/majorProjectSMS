package menus

// import (
// 	"context"
// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/handlers"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/fatih/color"
// )

// func ShowAdminInvoiceMenu(ctx context.Context, iHandler *handlers.InvoiceHandler) {
// 	for {
// 		color.Cyan("1." + string(constants.IssueInvoice))
// 		color.Cyan("2." + string(constants.SearchAInvoice))
// 		color.Cyan("3." + string(constants.ListInvoicesOfAYear))
// 		color.Cyan("4. Exit")
// 		for {
// 			ch := utils.ReadChoice()

// 			if ch == "" {
// 				break
// 			}
// 			switch ch {
// 			case "1":
// 				iHandler.IssueInvoice(ctx)
// 			case "2":
// 				iHandler.GetInvoiceByMonthAndYear()
// 			case "3":
// 				iHandler.GetInvoicesByYear()
// 			case "4":
// 				color.Red("Exit")
// 				return
// 			default:
// 				color.Red("Invalid choice, try again.")
// 			}
// 		}
// 	}
// }
