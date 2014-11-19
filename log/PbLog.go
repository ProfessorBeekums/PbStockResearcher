// Wrapper for log.Logger. The purpose behind this is to change where
// log messages go to by only changing this file.

package log

import (
	"log"
	"os"
)

var printLogger *log.Logger
var errorLogger *log.Logger

func getErrorLogger() *log.Logger {
	if errorLogger == nil {
		// TODO allow configuring an error file here?
		errorLogger = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile) 
	}

	return errorLogger
}

func getPrintLogger() *log.Logger {
	if printLogger == nil {
		printLogger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}

	return printLogger
}

func Println(v ...interface{}) {
	getPrintLogger().Println(v...)
}

func Fatal(v ...interface{}) {
	getErrorLogger().Fatal(v...)
}

func Error(v ...interface{}) {
	getErrorLogger().Println(v...)
}
