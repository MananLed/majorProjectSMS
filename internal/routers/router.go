package routers

import (
	"net/http"

	"github.com/MananLed/majorProjectSMS/internal/middleware"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/web_handlers"
)

func SetupRouter(userService service.UserService, serviceRequestService service.ServiceRequestService, noticeService service.NoticeService, societyService service.SocietyService, feedbackService service.FeedbackService, credentialService service.CredentialService, invoiceService service.InvoiceService) *http.ServeMux {
	r := http.NewServeMux()

	userHandler := web_handlers.NewUserHandler(&userService, &serviceRequestService)
	serviceRequestHandler := web_handlers.NewServiceRequestHandler(&serviceRequestService)
	noticeHandler := web_handlers.NewNoticeHandler(&noticeService)
	feedbackHandler := web_handlers.NewFeedbackHandler(&feedbackService)
	credentialHandler := web_handlers.NewCredentialHandler(&credentialService)
	societyHandler := web_handlers.NewSocietyHandler(&societyService)
	invoiceHandler := web_handlers.NewInvoiceHandler(&invoiceService)

	r.HandleFunc("POST /signup", http.HandlerFunc(userHandler.SignUp))  //Content in body
	r.HandleFunc("POST /login", http.HandlerFunc(userHandler.Login))  //content in body
	r.Handle("GET /profile", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.ViewProfile)))          // GET
	r.Handle("PATCH /profile/update", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.UpdateProfile))) // PATCH  //update details in body
	r.Handle("PATCH /profile/password", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.ChangePassword)))
	r.Handle("DELETE /profile", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.DeleteProfile))) // DELETE
	r.Handle("POST /officers", middleware.AuthMiddleWare(http.HandlerFunc(userHandler.CreateOfficer)))     // POST  //detail of officer in body

	r.Handle("POST /service", middleware.AuthMiddleWare(http.HandlerFunc(serviceRequestHandler.BookServiceRequest)))        // POST   //service details in body, type of service in query param
	r.Handle("PATCH /service/reschedule/{id}", middleware.AuthMiddleWare(http.HandlerFunc(serviceRequestHandler.RescheduleServiceRequest)))   // PATCH   //service new time slot in body
	r.Handle("DELETE /service/cancel/{id}", middleware.AuthMiddleWare(http.HandlerFunc(serviceRequestHandler.CancelServiceRequest)))           // DELETE  //doubt!!! id dynamic in path param
	r.Handle("GET /service", middleware.AuthMiddleWare(http.HandlerFunc(serviceRequestHandler.GetRequestsOfResident)))     // GET   // status in query param gets requests of a particular resident by status
	r.Handle("PATCH /service/approve/{id}", middleware.AuthMiddleWare(http.HandlerFunc(serviceRequestHandler.ApproveRequest)))  // PATCH   //explore path param for service id
	r.Handle("GET /service/time-slots", middleware.AuthMiddleWare(http.HandlerFunc(serviceRequestHandler.GetAvailableTimeSlots))) // GET
	r.Handle("GET /service/type-status", middleware.AuthMiddleWare(http.HandlerFunc(serviceRequestHandler.GetRequestsByServiceTypeAndStatus))) //GET
	//TODO: Some more method in service request regarding admin and officer to view approved and pending service requests


	r.Handle("GET /notices", middleware.AuthMiddleWare(http.HandlerFunc(noticeHandler.GetNotices)))                       // GET
	r.Handle("POST /notices/issue", middleware.AuthMiddleWare(http.HandlerFunc(noticeHandler.IssueNotice)))                // POST
	r.Handle("GET /notices/month-year", middleware.AuthMiddleWare(http.HandlerFunc(noticeHandler.GetNoticesByMonthYear))) // GET //month and year in query param

	r.Handle("POST /invoices/issue", middleware.AuthMiddleWare(http.HandlerFunc(invoiceHandler.IssueInvoice)))
	r.Handle("GET /invoices/month-year", middleware.AuthMiddleWare(http.HandlerFunc(invoiceHandler.GetInvoiceByMonthAndYear)))

	r.Handle("GET /society/residents", middleware.AuthMiddleWare(http.HandlerFunc(societyHandler.ViewResidents))) // GET
	r.Handle("GET /society/officers", middleware.AuthMiddleWare(http.HandlerFunc(societyHandler.ViewOfficers)))   // GET

	r.Handle("GET /feedbacks", middleware.AuthMiddleWare(http.HandlerFunc(feedbackHandler.GetFeedbacks))) // GET Sort it out while coding
	r.Handle("POST /feedbacks", middleware.AuthMiddleWare(http.HandlerFunc(feedbackHandler.GiveFeedback))) // POST   

	r.Handle("DELETE /credentials/officer", middleware.AuthMiddleWare(http.HandlerFunc(credentialHandler.DeleteOfficer))) // DELETE
	r.Handle("DELETE /credentials/resident", middleware.AuthMiddleWare(http.HandlerFunc(credentialHandler.DeleteResident))) // DELETE

	return r
}
