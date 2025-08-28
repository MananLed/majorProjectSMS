package menus

// import (

// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/handlers"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/fatih/color"
// )

// func ShowResidentNoticeMenu(nHandler *handlers.NoticeHandler) {
// 	for {
// 		color.Cyan("1." + string(constants.GetNoticePrompt))
// 		color.Cyan("2." + string(constants.GetNoticesByMonthYear))
// 		color.Cyan("3." + string(constants.GetNoticesByYear))
// 		color.Cyan("4. Exit")
// 		for {
// 			ch := utils.ReadChoice()
// 			if ch == "" {
// 				break
// 			}

// 			switch ch {
// 			case "1":
// 				nHandler.GetNotices()
// 			case "2":
// 				nHandler.GetNoticesByMonthYear()
// 			case "3":
// 				nHandler.GetNoticesByYear()
// 			case "4":
// 				color.Red("Exit")
// 				return
// 			default:
// 				color.Red("Invalid choice, try again.")
// 			}
// 		}
// 	}
// }
