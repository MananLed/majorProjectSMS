package menus

import (
	"context"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/handlers"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/fatih/color"
)

func ShowAdminOfficerMenu(ctx context.Context, socHandler *handlers.SocietyHandler, cHandler *handlers.CredentialHandler) {
	for {
		color.Cyan("1." + string(constants.ViewOfficersPrompt))
		color.Cyan("2." + string(constants.DeleteOfficerPrompt))
		color.Cyan("3. Exit")
		for {
			ch := utils.ReadChoice()
			if ch == "" {
				break
			}
			switch ch {
			case "1":
				socHandler.HandleViewOfficers(ctx)
			case "2":
				cHandler.DeleteOfficer(ctx)
			case "3":
				color.Red("Exit")
				return
			default:
				color.Red("Invalid choice, try again.")
			}
		}
	}
}
