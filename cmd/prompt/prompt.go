package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Stage prints a stage message
func Stage(msg string, a ...any) {
	fmt.Printf("[*] "+msg+"\n", a...)
}

// Action prints an action message
func Action(msg string, a ...any) {
	fmt.Printf("[+] "+msg+"\n", a...)
}

// Success prints a success message
func Success(msg string, a ...any) {
	fmt.Printf("[✓] "+msg+"\n", a...)
}

// Warn prints a warning message
func Warn(msg string, a ...any) {
	fmt.Printf("[!] "+msg+"\n", a...)
}

// Fail prints an error message and returns it as an error
func Fail(msg string, a ...any) error {
	fmt.Printf("[x] "+msg+"\n", a...)
	return fmt.Errorf(msg, a...)
}

// Ask prints a prompt message and returns user input
func Ask(msg string, a ...any) string {
	fmt.Printf("[?] "+msg, a...)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
