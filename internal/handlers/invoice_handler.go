package handlers

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/MananLed/majorProjectSMS/internal/model"
	"github.com/MananLed/majorProjectSMS/internal/service"
	"github.com/MananLed/majorProjectSMS/internal/utils"
	"github.com/MananLed/majorProjectSMS/pkg/logger"
	"github.com/fatih/color"
)

type InvoiceHandler struct {
	InvoiceService *service.InvoiceService
}

func NewInvoiceHandler(service *service.InvoiceService) *InvoiceHandler {
	return &InvoiceHandler{InvoiceService: service}
}

func (h *InvoiceHandler) IssueInvoice(ctx context.Context) {
	user, err := utils.GetUserFromContext(ctx)

	if user.Role != model.RoleAdmin {
		color.Red("not permitted to issue invoice")
		return
	}
	if err != nil {
		fmt.Print(color.RedString("error: "))
		fmt.Println(err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	var amount float64
	fmt.Print(color.YellowString("Enter the amount: "))
	fmt.Scanf("%f\n", &amount)

	now := time.Now()
	year := now.Year()
	month := now.Month()

	yearStr := fmt.Sprintf("%d", year)
	monthStr := month.String()

	err = h.InvoiceService.GenerateInvoice(amount, monthStr, yearStr)

	if err != nil {
		color.Red("Failed to issue invoice: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	} else {
		color.Green("Invoice issued successfully.")
	}
}

func (h *InvoiceHandler) GetInvoiceByMonthAndYear() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("Enter the month(MM): "))
	var monthindex int
	fmt.Scanf("%d\n", &monthindex)
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}

	for {
		if monthindex < 1 || monthindex > 12 {
			color.Red("invalid month, enter again")
		} else {
			break
		}
		fmt.Scanf("%d\n", &monthindex)
	}

	month := months[monthindex-1]
	month = strings.TrimRight(month, "\r\n")
	fmt.Print(color.YellowString("Enter the year(YYYY): "))

	year, _ := reader.ReadString('\n')
	year = strings.TrimRight(year, "\r\n")

	invoice, err := h.InvoiceService.GetInvoiceByMonthAndYear(month, year)

	if err != nil {
		fmt.Print(color.RedString("error: "))
		fmt.Println(err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.White(constants.InvoiceFormatPrompt, invoice.ID, invoice.Amount, invoice.Month, invoice.Year)
}

func (h *InvoiceHandler) GetInvoicesByYear() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("Enter the year(YYYY): "))
	year, _ := reader.ReadString('\n')
	year = strings.TrimRight(year, "\r\n")

	invoices, err := h.InvoiceService.GetInvoicesByYear(year)

	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	if len(invoices) > 0 {
		color.Green("Invoices:- ")
	}
	for _, invoice := range invoices {
		color.White(constants.InvoiceFormatPrompt, invoice.ID, invoice.Amount, invoice.Month, invoice.Year)
	}
}
