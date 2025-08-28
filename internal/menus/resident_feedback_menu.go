package menus

// import (
// 	"context"

// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/handlers"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/fatih/color"
// )

// func ShowResidentFeedbackMenu(ctx context.Context, fHandler *handlers.FeedbackHandler) {
// 	for {
// 		color.Cyan("1." + string(constants.GetFeedbackForResident))
// 		color.Cyan("2." + string(constants.GiveFeedbackPrompt))
// 		color.Cyan("3. Exit")
// 		for {
// 			ch := utils.ReadChoice()
// 			if ch == "" {
// 				break
// 			}

// 			switch ch {
// 			case "1":
// 				fHandler.GetFeebacksOfResident(ctx)
// 			case "2":
// 				fHandler.IssueFeedback(ctx)
// 			case "3":
// 				color.Red("Exit")
// 				return
// 			default:
// 				color.Red("Invalid choice, try again.")
// 			}
// 		}
// 	}
// }
