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

type NoticeHandler struct {
	NoticeService *service.NoticeService
}

func NewNoticeHandler(service *service.NoticeService) *NoticeHandler {
	return &NoticeHandler{NoticeService: service}
}

func (h *NoticeHandler) IssueNotice(ctx context.Context) {
	user, err := utils.GetUserFromContext(ctx)

	if user.Role == model.RoleResident {
		color.Red("not permitted to issue notice.")
		return
	}
	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("Enter notice content: "))
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	if content == "" {
		color.Red("Notice content cannot be empty.")
		return
	}

	now := time.Now()
	year := now.Year()
	month := now.Month()

	yearStr := fmt.Sprintf("%d", year)
	monthStr := month.String()

	err = h.NoticeService.IssueNotice(content, monthStr, yearStr)
	if err != nil {
		color.Red("Failed to issue notice: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
	} else {
		color.Green("Notice issued successfully.")
	}
}

func (h *NoticeHandler) GetNotices() {
	notices, err := h.NoticeService.GetNotices()
	if err != nil {
		color.Red("failed to retrieve notices: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	if len(notices) == 0 {
		color.Yellow("No notices found.")
		return
	}

	color.Cyan("===== All Notices =====\n\n")
	for _, notice := range notices {
		color.White(constants.NoticeFormatPrompt, notice.ID, notice.DateIssued, notice.Content)
	}
}

func (h *NoticeHandler) GetNoticesByMonthYear() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("Enter the month: "))
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

	notices, err := h.NoticeService.GetNoticesByMonthYear(month, year)
	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}
	fmt.Println(len(notices))
	color.Green("Notices Found:-")
	for _, notice := range notices {
		color.White(constants.NoticeFormatPrompt, notice.ID, notice.DateIssued, notice.Content)
	}
}

func (h *NoticeHandler) GetNoticesByYear() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.YellowString("Enter the year(YYYY): "))
	year, _ := reader.ReadString('\n')
	year = strings.TrimRight(year, "\r\n")
	notices, err := h.NoticeService.GetNoticesByYear(year)

	if err != nil {
		color.Red("error: %v", err)
		logger.LogToFile(fmt.Sprintf("error: %v", err))
		return
	}

	color.Green("Notices:- ")
	for _, notice := range notices {
		color.White(constants.NoticeFormatPrompt, notice.ID, notice.DateIssued, notice.Content)
	}
}
