package logging

import (
	"io"
	"log"
	"os"
)

const logFlag = log.Ldate | log.Ltime | log.Llongfile

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {

	file, err := os.OpenFile("errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error opening file")
	}

	Info = log.New(os.Stdout, "INFO: ", logFlag)
	Warning = log.New(os.Stdout, "WARNING: ", logFlag)
	Error = log.New(io.MultiWriter(os.Stderr, file), "ERROR: ", logFlag)
}
