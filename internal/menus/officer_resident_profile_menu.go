package menus

// import (
// 	"context"

// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/handlers"
// 	"github.com/MananLed/majorProjectSMS/internal/model"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/fatih/color"
// )

// func ShowOfficerResidentProfileMenu(ctx context.Context, user *model.User, uHandler *handlers.UserHandler) bool {
// 	for {
// 		color.Cyan("1." + string(constants.UpdateProfilePrompt))
// 		color.Cyan("2." + string(constants.ChangePasswordPrompt))
// 		color.Cyan("3." + string(constants.ViewProfilePrompt))
// 		color.Red("4." + string(constants.DeleteProfilePrompt))
// 		color.Cyan("5." + "Exit")
// 		for {
// 			ch := utils.ReadChoice()
// 			if ch == "" {
// 				break
// 			}
// 			switch ch {
// 			case "1":
// 				uHandler.UpdateProfile(user)
// 			case "2":
// 				uHandler.ChangePassword(user)
// 			case "3":
// 				uHandler.ViewProfile(user)
// 			case "4":
// 				uHandler.DeleteProfile(ctx)
// 				return true
// 			case "5":
// 				color.Red("Exit")
// 				return false
// 			default:
// 				color.Red("Invalid choice, try again.")
// 			}
// 		}
// 	}
// }
