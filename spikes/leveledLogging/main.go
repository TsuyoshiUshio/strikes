package main

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

const (
	EnvLog     = "ST_LOG"
	EnvLogFile = "ST_LOG_PATH"
)

func main() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "ERROR"},
		MinLevel: logutils.LogLevel(os.Getenv(EnvLog)),
		Writer:   os.Stderr,
	}
	log.SetOutput(filter)
	log.SetPrefix("[Strikes] ")
	log.Println("[DEBUG] Now debugging...")
	log.Println("[INFO] Echo useful information...")
	log.Println("[ERROR] Error happens!...")
	log.Println("Normal logging")
	log.Fatalln("[DEBUG] Fatail!")

}
