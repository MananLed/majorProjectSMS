package utils

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/fatih/color"
)

func PromptRequired(label string, reader *bufio.Reader) string {
	for {
		fmt.Print(color.YellowString(label) + ": ")
		input, _ := reader.ReadString('\n')
		trimmed := strings.TrimSpace(input)
		trimmed = strings.TrimRight(trimmed, "\r\n")
		if trimmed == "" {
			color.Red("%s is compulsory", label)
			continue
		}
		return trimmed
	}
}

func ValidateMobileNumber(mobile string) bool {
	mobile = strings.TrimSpace(mobile)

	pattern := `^[6-9][0-9]{9}$`

	re := regexp.MustCompile(pattern)
	return re.MatchString(mobile)
}

func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		color.Red("Invalid email format.")
		return false
	}
	return true
}

func ValidatePassword(password string) bool {
	var hasLower, hasDigit, hasSpecial bool

	if len(password) < 12 {
		color.Red("Password must be at least 12 characters long.")
		return false
	}

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower || !hasDigit || !hasSpecial {
		color.Red("Password must contain at least one lowercase letter, one digit, and one special character.")
		return false
	}

	return true
}

func ValidateFlatNumber(flat string) bool {
	pattern := `^[0-8]0[1-4]$`
	matched, _ := regexp.MatchString(pattern, flat)
	return matched
}