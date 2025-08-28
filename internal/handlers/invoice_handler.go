package handlers

// import (

// 	"context"
// 	"fmt"

// 	"time"

// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/model"
// 	"github.com/MananLed/majorProjectSMS/internal/service"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/MananLed/majorProjectSMS/pkg/logger"
// 	"github.com/fatih/color"
// )

// type InvoiceHandler struct {
// 	InvoiceService service.InvoiceServiceInterface
// }

// func NewInvoiceHandler(service service.InvoiceServiceInterface) *InvoiceHandler {
// 	return &InvoiceHandler{InvoiceService: service}
// }


// func (h *InvoiceHandler) IssueInvoice(ctx context.Context) {
// 	user, err := utils.GetUserFromContext(ctx)
// 	if err != nil {
// 		fmt.Print(color.RedString("error: "))
// 		fmt.Println(err)
// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
// 		return
// 	}

// 	if user.Role != model.RoleAdmin {
// 		color.Red("not permitted to issue invoice")
// 		return
// 	}

// 	var amount float64
// 	fmt.Print(color.YellowString("Enter the amount: "))
// 	fmt.Scanf("%f\n", &amount)

// 	now := time.Now()
// 	year := now.Year()    
// 	month := now.Month()    

// 	err = h.InvoiceService.GenerateInvoice(amount, month, year)
// 	if err != nil {
// 		color.Red("Failed to issue invoice: %v", err)
// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
// 	} else {
// 		color.Green("Invoice issued successfully.")
// 	}
// }


// func (h *InvoiceHandler) GetInvoiceByMonthAndYear() {
// 	fmt.Print(color.YellowString("Enter the month (1-12): "))
// 	var monthIndex int
// 	fmt.Scanf("%d\n", &monthIndex)

// 	for monthIndex < 1 || monthIndex > 12 {
// 		color.Red("Invalid month, enter again: ")
// 		fmt.Scanf("%d\n", &monthIndex)
// 	}
// 	month := time.Month(monthIndex)

// 	var year int
// 	fmt.Print(color.YellowString("Enter the year (YYYY): "))
// 	fmt.Scanf("%d\n", &year)

// 	invoice, err := h.InvoiceService.GetInvoiceByMonthAndYear(month, year)
// 	if err != nil {
// 		fmt.Print(color.RedString("error: "))
// 		fmt.Println(err)
// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
// 		return
// 	}

// 	color.White(constants.InvoiceFormatPrompt, invoice.ID, invoice.Amount, invoice.Month.String(), invoice.Year)
// }

// func (h *InvoiceHandler) GetInvoicesByYear() {
// 	var year int
// 	fmt.Print(color.YellowString("Enter the year (YYYY): "))
// 	fmt.Scanf("%d\n", &year)

// 	invoices, err := h.InvoiceService.GetInvoicesByYear(year)
// 	if err != nil {
// 		color.Red("error: %v", err)
// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
// 		return
// 	}

// 	if len(invoices) > 0 {
// 		color.Green("Invoices:- ")
// 	} else {
// 		color.Yellow("No invoices found for %d", year)
// 	}

// 	for _, invoice := range invoices {
// 		color.White(constants.InvoiceFormatPrompt, invoice.ID, invoice.Amount, invoice.Month.String(), invoice.Year)
// 	}
// }