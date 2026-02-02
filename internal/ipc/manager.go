package ipc

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"github.com/unf6/nucleus/internal/config"
)

type ipcTargets map[string]map[string][]string

var qsHasAnyDisplay = sync.OnceValue(func() bool {
	out, err := exec.Command("qs", "ipc", "--help").Output()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), "--any-display")
})

func ParseTargetsFromIPCShowOutput(output string) ipcTargets {
	targets := make(ipcTargets)
	var currentTarget string
	for _, line := range strings.Split(output, "\n") {
		if after, ok := strings.CutPrefix(line, "target "); ok {
			currentTarget = strings.TrimSpace(after)
			targets[currentTarget] = make(map[string][]string)
		}
		if strings.HasPrefix(line, "  function") && currentTarget != "" {
			argsList := []string{}
			currentFunc := strings.TrimPrefix(line, "  function ")
			funcDef := strings.SplitN(currentFunc, "(", 2)
			if len(funcDef) < 2 {
				continue
			}
			argList := strings.SplitN(funcDef[1], ")", 2)[0]
			args := strings.Split(argList, ",")
			if len(args) > 0 && strings.TrimSpace(args[0]) != "" {
				argsList = append(argsList, funcDef[0])
				for _, arg := range args {
					argName := strings.SplitN(strings.TrimSpace(arg), ":", 2)[0]
					argsList = append(argsList, argName)
				}
				targets[currentTarget][funcDef[0]] = argsList
			} else {
				targets[currentTarget][funcDef[0]] = make([]string, 0)
			}
		}
	}
	return targets
}

func GetShellIPCCompletions(args []string, _ string) []string {
	shellFile, err := config.GetShellFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting shell file path: %v\n", err)
		return nil
	}

	cmdArgs := []string{"ipc"}
	if qsHasAnyDisplay() {
		cmdArgs = append(cmdArgs, "--any-display")
	}
	cmdArgs = append(cmdArgs, "-p", shellFile, "show")
	cmd := exec.Command("qs", cmdArgs...)
	var targets ipcTargets

	if output, err := cmd.Output(); err == nil {
		targets = ParseTargetsFromIPCShowOutput(string(output))
	} else {
		fmt.Fprintf(os.Stderr, "Error getting IPC show output for completions: %v\n", err)
		return nil
	}

	if len(args) > 0 && args[0] == "call" {
		args = args[1:]
	}

	if len(args) == 0 {
		targetNames := make([]string, 0)
		targetNames = append(targetNames, "call")
		for k := range targets {
			targetNames = append(targetNames, k)
		}
		return targetNames
	}
	if len(args) == 1 {
		if targetFuncs, ok := targets[args[0]]; ok {
			funcNames := make([]string, 0)
			for k := range targetFuncs {
				funcNames = append(funcNames, k)
			}
			return funcNames
		}
		return nil
	}
	if targetFuncs, ok := targets[args[0]]; ok {
		if funcArgs, ok := targetFuncs[args[1]]; ok && len(args) <= len(funcArgs) {
			if len(funcArgs) >= len(args) {
				return []string{fmt.Sprintf("[%s]", funcArgs[len(args)-1])}
			}
		}
	}

	return nil
}

func RunShellIPCCommand(args []string) {
	shellFile, err := config.GetShellFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting shell file path: %v\n", err)
		os.Exit(1)
	}

	if len(args) == 0 {
		PrintIPCHelp(shellFile)
		return
	}

	if args[0] != "call" {
		args = append([]string{"call"}, args...)
	}

	cmdArgs := []string{"ipc"}
	if qsHasAnyDisplay() {
		cmdArgs = append(cmdArgs, "--any-display")
	}
	cmdArgs = append(cmdArgs, "-p", shellFile)
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("qs", cmdArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running IPC command: %v\n", err)
		os.Exit(1)
	}
}

func PrintIPCHelp(shellFile string) {
	fmt.Println("Usage: nucleus ipc <target> <function> [args...]")
	fmt.Println()

	cmdArgs := []string{"ipc"}
	if qsHasAnyDisplay() {
		cmdArgs = append(cmdArgs, "--any-display")
	}
	cmdArgs = append(cmdArgs, "-p", shellFile, "show")
	cmd := exec.Command("qs", cmdArgs...)

	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Could not retrieve available IPC targets (is shell running?)")
		return
	}

	targets := ParseTargetsFromIPCShowOutput(string(output))
	if len(targets) == 0 {
		fmt.Println("No IPC targets available")
		return
	}

	fmt.Println("Targets:")

	targetNames := make([]string, 0, len(targets))
	for name := range targets {
		targetNames = append(targetNames, name)
	}
	slices.Sort(targetNames)

	for _, targetName := range targetNames {
		funcs := targets[targetName]
		funcNames := make([]string, 0, len(funcs))
		for fn := range funcs {
			funcNames = append(funcNames, fn)
		}
		slices.Sort(funcNames)
		fmt.Printf("  %-16s %s\n", targetName, strings.Join(funcNames, ", "))
	}
}
