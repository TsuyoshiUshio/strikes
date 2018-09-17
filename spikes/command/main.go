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

}

func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "--help")
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
