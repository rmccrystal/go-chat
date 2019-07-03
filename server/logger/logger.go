package logger

import (
	"fmt"
	"time"
)

// Logs a message with prefix, for example if the prefix was ERROR
// the output would be [time] - ERROR - (text)
func LogPrefixf(format string, prefix string, args ...interface{}) {
	logPrefix := fmt.Sprintf("[%s] » %s » ", time.Now().Format("2006-01-02 15:04:05"), prefix)
	fmt.Printf(logPrefix+format, args...)
}

func Logf(format string, args ...interface{}) {
	prefix := fmt.Sprintf("[%s] » ", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Printf(prefix+format, args...)
}

func Errorf(format string, args ...interface{}) {
	LogPrefixf(format, "ERROR", args...)
}
