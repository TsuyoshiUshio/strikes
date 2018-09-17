package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	command := os.Args[1]
	// What do I want to achieve in this spike
	// 1. How to check the external command
	fmt.Println("Available: " + command + " " + strconv.FormatBool(isCommandAvailable(command)))
	// 2. How to execute the extenarl command with flexible command
	fmt.Println("start executing terraform init")
	commandAndArgs := []string{
		"init",
	}
	cmd := exec.Command("terraform", commandAndArgs...) // ...enable us to pass them slice
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run() // Run for Wait for the execution.
	if err != nil {
		panic(err)
	}
	fmt.Println("finished.")
}

func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "--help") // terraform works. it is not for everyCommand. Make it simple. it should depend on provider.
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
