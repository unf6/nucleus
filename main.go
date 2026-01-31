package main

import (
        "github.com/unf6/nucleus/cmd"
        "github.com/charmbracelet/log"

        "os"
)

func main() {
        if err := cmd.Execute(); err != nil {
                log.Error("Command Execution Failed: ", err)
                os.Exit(1)
        }
} 
