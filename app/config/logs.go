package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func InitLogger() error {
	logFolder := "./logs"
	timestamp := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("%s/server_%s.log", logFolder, timestamp)

	var err error
	logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	Logger = log.New(multiWriter, "FORUM: ", log.Ldate|log.Ltime|log.Lshortfile)
	Logger.Println("Logger created successfully")
	return nil
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
