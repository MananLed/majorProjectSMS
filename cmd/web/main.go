package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/config"
	"github.com/MananLed/majorProjectSMS/internal/db"
	"github.com/MananLed/majorProjectSMS/internal/middleware"
	"github.com/MananLed/majorProjectSMS/internal/repository"
	"github.com/MananLed/majorProjectSMS/internal/routers"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
)

func main() {

	database, err := db.InitDB()
	if err != nil {
		log.Fatal(err)
		logger.LogToFile(fmt.Sprintf("Error: %v", err))
	}
	defer database.Close()
	if err := db.RunInitialSetup(database); err != nil {
		log.Fatal(err)
		logger.LogToFile(fmt.Sprintf("Error: %v", err))
	}
	userRepo := repository.NewUserRepository(database)
	serviceRequestRepo := repository.NewServiceRequestRepository(database)
	noticeRepo := repository.NewNoticeRepository(database)
	societyRepo := repository.NewSocietyRepository(database)
	feedbackRepo := repository.NewFeedbackRepository(database)
	credentialRepo := repository.NewCredentialRepository(database)
	invoiceRepo := repository.NewInvoiceRepository(database)

	userService := service.NewUserService(userRepo)
	serviceRequestService := service.NewServiceRequestService(serviceRequestRepo)
	noticeService := service.NewNoticeService(noticeRepo)
	societyService := service.NewSocietyService(societyRepo)
	feedbackService := service.NewFeedbackService(feedbackRepo)
	credentialService := service.NewCredentialService(credentialRepo)
	invoiceService := service.NewInvoiceService(invoiceRepo)

	router := routers.SetupRouter(*userService, *serviceRequestService, *noticeService, *societyService, *feedbackService, *credentialService, *invoiceService)

	handler := middleware.CorsMiddleWare(
		middleware.LoggingMiddleWare(router),
	)

	log.Println("Server stating on the Port 8080...")
	fmt.Println("Server stating on the Port 8080...")

	err = http.ListenAndServe(config.PortForAPI, handler)
	if err != nil {
		log.Fatal(err)
	}

	// mux.HandleFunc("/signup", userHandler.SignUp)  //Content in body
	// mux.HandleFunc("/login", userHandler.Login)  //content in body
	// mux.HandleFunc("/profile", userHandler.ViewProfile)          // GET
	// mux.HandleFunc("/profile/update", userHandler.UpdateProfile) // PATCH  //update details in body
	// mux.HandleFunc("/profile/password", userHandler.ChangePassword)
	// mux.HandleFunc("/profile/delete", userHandler.DeleteProfile) // DELETE
	// mux.HandleFunc("/officers", userHandler.CreateOfficer)       // POST  //detail of officer in body

	// mux.HandleFunc("/service/book", serviceRequestHandler.BookServiceRequest())        // POST   //service details in body, type of service in query param
	// mux.HandleFunc("/service/reschedule", serviceRequestHandler.RescheduleServiceRequest())   // PATCH   //service new time slot in body
	// mux.HandleFunc("/service/cancel", serviceRequestHandler.CancelServiceRequest())           // DELETE  //doubt!!! id dynamic in path param
	// mux.HandleFunc("/service/status", serviceRequestHandler.GetRequestsByStatus())     // GET   // status in query param
	// mux.HandleFunc("/service/approve", serviceRequestHandler.ApproveRequest())  // PATCH   //explore path param for service id
	// mux.HandleFunc("/service/time-slots", serviceRequestHandler.GetAvailableTimeSlots()) // GET

	// mux.HandleFunc("/notices", noticeHandler.GetNotices())                       // GET
	// mux.HandleFunc("/notices/issue", noticeHandler.IssueNotice())                // POST
	// mux.HandleFunc("/notices/month-year", noticeHandler.GetNoticesByMonthYear) // GET //month and year in query param
	// mux.HandleFunc("/notices/year", noticeHandler.GetNoticesByYear)            // GET //year in query param

	// mux.HandleFunc("/society/residents", societyHandler.HandleViewResidents) // GET
	// mux.HandleFunc("/society/officers", societyHandler.HandleViewOfficers)   // GET

	// mux.HandleFunc("/feedbacks", feedbackHandler.GetFeedbacks) // GET
	// mux.HandleFunc("/feedbacks", feedbackHandler.GiveFeedback) // POST   Sort it out while coding

	// mux.HandleFunc("/credentials/delete", credentialHandler.DeleteCredential) // DELETE

}
