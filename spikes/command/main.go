package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	stdout, err := cmd.StdoutPipe()                     // piping test
	if err != nil {
		log.Fatal(err)
	}
	var outBuf bytes.Buffer
	StdOutMulti := io.MultiWriter(&outBuf, os.Stdout)
	cmd.Stdout = StdOutMulti

	var buf bytes.Buffer
	multiWriter := io.MultiWriter(&buf, os.Stderr)
	cmd.Stderr = multiWriter
	err = cmd.Start() // Run for Wait for the execution.
	if err != nil {
		panic(err)
	}
	fmt.Println("-------poped output")
	result, _ := ioutil.ReadAll(stdout)
	fmt.Println(string(result)) // works

	if err = cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("finished.")

	fmt.Println("------- err output ")
	fmt.Println(buf.String())

	fmt.Println("-------stdout")
	fmt.Println(outBuf.String())

}

func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "--help") // terraform works. it is not for everyCommand. Make it simple. it should depend on provider.
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
