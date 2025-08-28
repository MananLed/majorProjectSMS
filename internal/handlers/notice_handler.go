package handlers

// import (
// 	"bufio"
// 	"context"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"

// 	"github.com/MananLed/majorProjectSMS/constants"
// 	"github.com/MananLed/majorProjectSMS/internal/model"
// 	"github.com/MananLed/majorProjectSMS/internal/service"
// 	"github.com/MananLed/majorProjectSMS/internal/utils"
// 	"github.com/MananLed/majorProjectSMS/pkg/logger"
// 	"github.com/fatih/color"
// )

// type NoticeHandler struct {
// 	NoticeService service.NoticeServiceInterface
// }

// func NewNoticeHandler(service service.NoticeServiceInterface) *NoticeHandler {
// 	return &NoticeHandler{NoticeService: service}
// }

// func (h *NoticeHandler) IssueNotice(ctx context.Context) {
// 	user, err := utils.GetUserFromContext(ctx)

// 	if err != nil {
// 		color.Red("error: %v", err)
// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
// 		return
// 	}
// 	if user.Role == model.RoleResident {
// 		color.Red("not permitted to issue notice.")
// 		return
// 	}
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print(color.YellowString("Enter notice content: "))
// 	content, _ := reader.ReadString('\n')
// 	content = strings.TrimSpace(content)

// 	if content == "" {
// 		color.Red("Notice content cannot be empty.")
// 		return
// 	}

// 	now := time.Now()


// 	err = h.NoticeService.IssueNotice(content, now.Month(), now.Year())
// 	if err != nil {
// 		color.Red("Failed to issue notice: %v", err)
// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
// 	} else {
// 		color.Green("Notice issued successfully.")
// 	}
// }

// func (h *NoticeHandler) GetNotices() {
// 	notices, err := h.NoticeService.GetNotices()
// 	if err != nil {
// 		color.Red("failed to retrieve notices: %v", err)
// 		logger.LogToFile(fmt.Sprintf("error: %v", err))
// 		return
// 	}

// 	if len(notices) == 0 {
// 		color.Yellow("No notices found.")
// 		return
// 	}

// 	color.Cyan("===== All Notices =====\n\n")
// 	for _, notice := range notices {
// 		dateStr := notice.DateIssued.Format("02-Jan-2006")
// 		color.White(constants.NoticeFormatPrompt, notice.ID, dateStr, notice.Content)
// 	}
// }

// func (h *NoticeHandler) GetNoticesByMonthYear() {
// 	reader := bufio.NewReader(os.Stdin)

// 	var monthIndex int
// 	for {
// 		fmt.Print(color.YellowString("Enter the month (1-12): "))
// 		_, err := fmt.Scanf("%d\n", &monthIndex)
// 		if err == nil && monthIndex >= 1 && monthIndex <= 12 {
// 			break
// 		}
// 		color.Red("Invalid month. Please enter a number between 1 and 12.")
// 	}

// 	fmt.Print(color.YellowString("Enter the year (YYYY): "))
// 	yearStr, _ := reader.ReadString('\n')
// 	yearStr = strings.TrimSpace(yearStr)

// 	year, err := strconv.Atoi(yearStr)
// 	if err != nil || year < 1 {
// 		color.Red("Invalid year.")
// 		return
// 	}

// 	notices, err := h.NoticeService.GetNoticesByMonthYear(time.Month(monthIndex), year)
// 	if err != nil {
// 		color.Red("Error retrieving notices: %v", err)
// 		logger.LogToFile(fmt.Sprintf("Error retrieving notices: %v", err))
// 		return
// 	}

// 	if len(notices) == 0 {
// 		color.Yellow("No notices found for %s %d.", time.Month(monthIndex), year)
// 		return
// 	}

// 	color.Green("Notices Found:")
// 	for _, notice := range notices {
// 		dateStr := notice.DateIssued.Format("02-Jan-2006")
// 		color.White(constants.NoticeFormatPrompt, notice.ID, dateStr, notice.Content)
// 	}
// }


// func (h *NoticeHandler) GetNoticesByYear() {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print(color.YellowString("Enter the year (YYYY): "))
// 	yearStr, _ := reader.ReadString('\n')
// 	yearStr = strings.TrimSpace(yearStr)

// 	year, err := strconv.Atoi(yearStr)
// 	if err != nil || year < 1 {
// 		color.Red("Invalid year.")
// 		return
// 	}

// 	notices, err := h.NoticeService.GetNoticesByYear(year)
// 	if err != nil {
// 		color.Red("Error retrieving notices: %v", err)
// 		logger.LogToFile(fmt.Sprintf("Error retrieving notices: %v", err))
// 		return
// 	}

// 	if len(notices) == 0 {
// 		color.Yellow("No notices found for year %d.", year)
// 		return
// 	}

// 	color.Green("Notices:")
// 	for _, notice := range notices {
// 		dateStr := notice.DateIssued.Format("02-Jan-2006")
// 		color.White(constants.NoticeFormatPrompt, notice.ID, dateStr, notice.Content)
// 	}
// }