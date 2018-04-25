package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Confirm wait for user response returning
// `true` if the answer was y or yes
func Confirm(message string) bool {
	fmt.Printf("%s [y/N] ", message)

	reader := bufio.NewReader(os.Stdin)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}
