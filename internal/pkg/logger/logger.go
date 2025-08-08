package logger

import (
	"log"
	"os"
)

// ROVERT
// infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
// errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

var (
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DebugLog *log.Logger
)

func Init() {
	InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	ErrorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLog = log.New(os.Stdout, "DEBUG\t", log.Ldate|log.Ltime|log.Lshortfile) // Oportunity
}
