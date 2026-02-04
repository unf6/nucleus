package prompt

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Stage(msg string, a ...any) {
	fmt.Printf("[*] "+msg+"\n", a...)
}

func Action(msg string, a ...any) {
	fmt.Printf("[+] "+msg+"\n", a...)
}

func Success(msg string, a ...any) {
	fmt.Printf("[✓] "+msg+"\n", a...)
}

func Warn(msg string, a ...any) {
	fmt.Printf("[!] "+msg+"\n", a...)
}

func Fail(msg string, a ...any) error {
	fmt.Printf("[x] "+msg+"\n", a...)
	return fmt.Errorf(msg, a...)
}

func Ask(msg string, a ...any) string {
	fmt.Printf("[?] "+msg, a...)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
