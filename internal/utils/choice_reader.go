package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/MananLed/majorProjectSMS/constants"
	"github.com/fatih/color"
)

func ReadChoice() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(color.BlueString(string(constants.ChoicePrompt)))
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	return choice
}
