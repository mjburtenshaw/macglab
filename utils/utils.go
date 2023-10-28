package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func AskBinaryQuestion(question string) (response string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	response, _ = reader.ReadString('\n')
	response = strings.TrimSpace(response)
	return response
}
