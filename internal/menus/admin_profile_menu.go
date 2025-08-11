package menus

import (
	"context"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/handlers"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/fatih/color"
)

func ShowAdminProfileMenu(ctx context.Context, uHandler *handlers.UserHandler, user *model.User) {
	for {
		color.Cyan("1." + string(constants.UpdateProfilePrompt))
		color.Cyan("2." + string(constants.ChangePasswordPrompt))
		color.Cyan("3." + string(constants.ViewProfilePrompt))
		color.Cyan("4. Exit")
		for {
			ch := utils.ReadChoice()
			if ch == "" {
				break
			}
			switch ch {
			case "1":
				uHandler.UpdateProfile(user)
			case "2":
				uHandler.ChangePassword(ctx)
			case "3":
				uHandler.ViewProfile(user)
			case "4":
				color.Red("Exit")
				return
			default:
				color.Red("Invalid choice, try again.")
			}
		}
	}
}
