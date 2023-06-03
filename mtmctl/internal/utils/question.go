package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskUser(question string) bool {
	fmt.Printf("%s (Y/N): ", question)

	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))

	return answer == "yes" || answer == "y"
}
