package helpers

import (
	"log"
	"os"

	"github.com/hashicorp/logutils"
)

const (
	EnvLog     = "ST_LOG"
	EnvLogFile = "ST_LOG_PATH"
)

func SetUpLogger() {
	value, isPresent := os.LookupEnv(EnvLog)
	var minLevel string
	if !isPresent {
		minLevel = "INFO"
	} else {
		if value == "DEBUG" || value == "INFO" || value == "ERROR" {
			minLevel = value
		} else {
			log.Fatalf("EnvironmentVariables ST_LOG has wrong value %v. It should be DEBUG, INFO, or ERROR", value)
		}
	}

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "INFO", "ERROR"},
		MinLevel: logutils.LogLevel(minLevel),
		Writer:   os.Stderr,
	}

	log.SetOutput(filter)
	log.SetPrefix("Strikes: ")
}
