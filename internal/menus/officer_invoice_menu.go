package menus

// import (

// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/handlers"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/fatih/color"
// )

// func ShowOfficerInvoiceMenu(iHandler *handlers.InvoiceHandler) {
// 	for {
// 		color.Cyan("1." + string(constants.SearchAInvoice))
// 		color.Cyan("2." + string(constants.ListInvoicesOfAYear))
// 		color.Cyan("3." + "Exit")
// 		for {
// 			ch := utils.ReadChoice()
// 			if ch == "" {
// 				break
// 			}
// 			switch ch {
// 			case "1":
// 				iHandler.GetInvoiceByMonthAndYear()
// 			case "2":
// 				iHandler.GetInvoicesByYear()
// 			case "3":
// 				color.Red("Exit")
// 				return
// 			default:
// 				color.Red("Invalid choice, try again.")
// 			}
// 		}
// 	}
// }
