package utils

import (
	"log"
	"os"
	"producer/models"
)

var (
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	fatalLogger = log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile)
)

func Info(msg string) {
	if models.AppStateVar.Debug {
		infoLogger.Println(msg)
	}
}

func Error(msg string) {
	if models.AppStateVar.Debug {
		errorLogger.Println(msg)
	}
}

func Fatal(msg string) {
	if models.AppStateVar.Debug {
		fatalLogger.Println(msg)
	}
}