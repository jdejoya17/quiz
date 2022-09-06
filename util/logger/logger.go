package logger

import (
	"io/fs"
	"log"
	"os"
)

// define logger
var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {
	// initialize logger
	LOG_NAME := "quiz.log"
	LOG_FLAGS := os.O_APPEND | os.O_CREATE | os.O_WRONLY
	var LOG_PERM fs.FileMode = 0666
	LOGGER_FLAGS := log.Ldate | log.Lmicroseconds | log.Llongfile | log.Lmsgprefix

	f, err := os.OpenFile(LOG_NAME, LOG_FLAGS, LOG_PERM)
	if err != nil {
		log.Fatal(err)
	}

	Info = log.New(f, "INFO: ", LOGGER_FLAGS)
	Warning = log.New(f, "WARNING: ", LOGGER_FLAGS)
	Error = log.New(f, "ERROR: ", LOGGER_FLAGS)
}
