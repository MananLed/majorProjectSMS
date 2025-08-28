package menus

// import (
// 	"context"
// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/handlers"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/fatih/color"
// )

// func ShowAdminNoticeMenu(ctx context.Context, nHandler *handlers.NoticeHandler) {
// 	for {
// 		color.Cyan("1." + string(constants.IssueNoticePrompt))
// 		color.Cyan("2." + string(constants.GetNoticePrompt))
// 		color.Cyan("3." + string(constants.GetNoticesByMonthYear))
// 		color.Cyan("4." + string(constants.GetNoticesByYear))
// 		color.Cyan("5. Exit")
// 		for {
// 			ch := utils.ReadChoice()
// 			if ch == "" {
// 				break
// 			}
// 			switch ch {
// 			case "1":
// 				nHandler.IssueNotice(ctx)
// 			case "2":
// 				nHandler.GetNotices()
// 			case "3":
// 				nHandler.GetNoticesByMonthYear()
// 			case "4":
// 				nHandler.GetNoticesByYear()
// 			case "5":
// 				color.Red("Exit")
// 				return
// 			default:
// 				color.Red("Invalid choice, try again.")
// 			}
// 		}
// 	}
// }
