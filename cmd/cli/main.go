package main

// import (
// 	"bufio"
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strings"

// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/db"
// 	"github.com/MananLed/majorProjectSMS/internal/handlers"
// 	"github.com/MananLed/majorProjectSMS/internal/model"
// 	"github.com/MananLed/majorProjectSMS/internal/repository"
// 	"github.com/MananLed/majorProjectSMS/internal/service"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/MananLed/majorProjectSMS/pkg/logger"

// 	"github.com/common-nighthawk/go-figure"
// 	"github.com/fatih/color"
// )

// func main() {

// 	database, err := db.InitDB()
// 	if err != nil {
// 		log.Fatal(err)
// 		logger.LogToFile(fmt.Sprintf("Error: %v", err))
// 	}
// 	defer database.Close()
// 	if err := db.RunInitialSetup(database); err != nil {
// 		log.Fatal(err)
// 		logger.LogToFile(fmt.Sprintf("Error: %v", err))
// 	}

// 	societyRepo := repository.NewSocietyRepository(database)
// 	societyService := service.NewSocietyService(societyRepo)
// 	societyHandler := handlers.NewSocietyHandler(societyService)

// 	credentialRepo := repository.NewCredentialRepository(database)
// 	credentialService := service.NewCredentialService(credentialRepo)
// 	credentialHandler := handlers.NewCredentialHandler(credentialService)

// 	noticeRepo := repository.NewNoticeRepository(database)
// 	noticeService := service.NewNoticeService(noticeRepo)
// 	noticeHandler := handlers.NewNoticeHandler(noticeService)

// 	serviceRequestRepo := repository.NewServiceRequestRepository(database)
// 	serviceRequestService := service.NewServiceRequestService(serviceRequestRepo)
// 	serviceRequestHandler := handlers.NewServiceRequestHandler(serviceRequestService)

// 	feedbackRepo := repository.NewFeedbackRepository(database)
// 	feedbackService := service.NewFeedbackService(feedbackRepo)
// 	feedbackHandler := handlers.NewFeedbackHandler(feedbackService)

// 	invoiceRepo := repository.NewInvoiceRepository(database)
// 	invoiceService := service.NewInvoiceService(invoiceRepo)
// 	invoiceHandler := handlers.NewInvoiceHandler(invoiceService)

// 	userRepo := repository.NewUserRepository(database)
// 	userService := service.NewUserService(userRepo)
// 	userHandler := handlers.NewUserHandler(userService, serviceRequestService)

// 	reader := bufio.NewReader(os.Stdin)

// 	for {
// 		myFigure := figure.NewColorFigure("UpKeepz", "", "green", false)
// 		fmt.Println(constants.AppEmogiPrompt)
// 		myFigure.Print()
// 		fmt.Println(constants.AppEmogiPrompt)

// 		color.Cyan("1." + string(constants.SignUpPrompt))
// 		color.Cyan("2." + string(constants.LoginPrompt))
// 		color.Cyan("3." + string(constants.ExitPrompt))

// 		fmt.Print(color.BlueString(string(constants.ChoiceMainPrompt)))

// 		input, _ := reader.ReadString('\n')
// 		choice := strings.TrimSpace(input)

// 		switch choice {
// 		case "1":
// 			userHandler.SignUp()
// 		case "2":
// 			user := userHandler.Login()
// 			if user == nil {
// 				continue
// 			}

// 			ctx := context.Background()
// 			ctx = context.WithValue(ctx, utils.UserIDKey, user.ID)
// 			ctx = context.WithValue(ctx, utils.UserRoleKey, user.Role)
// 			ctx = context.WithValue(ctx, utils.UserEmailKey, user.Email)
// 			ctx = context.WithValue(ctx, utils.UserFlatKey, user.Flat)

// 			switch user.Role {
// 			case model.RoleAdmin:
// 				ShowAdminDashboard(ctx, user, userHandler, societyHandler, credentialHandler, noticeHandler, feedbackHandler, invoiceHandler, serviceRequestHandler)
// 			case model.RoleOfficer:
// 				ShowOfficerDashboard(ctx, user, userHandler, serviceRequestHandler, noticeHandler, feedbackHandler, invoiceHandler)
// 			case model.RoleResident:
// 				ShowResidentDashboard(ctx, user, userHandler, serviceRequestHandler, noticeHandler, feedbackHandler, invoiceHandler)
// 			default:
// 				color.Green("Logged in as %s", user.Role)
// 			}
// 		case "3":
// 			color.Red("Exit")
// 			return
// 		default:
// 			color.Red("Invalid choice. Please try again.")
// 		}
// 	}
// }
