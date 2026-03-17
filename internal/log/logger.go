package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

var DebugMode = false

func LogInfo(message string) error {
	_, err := fmt.Fprintln(os.Stdout, addTime(message))
	return err
}

func LogDebug(message string) error {
	if !DebugMode {
		return nil
	}

	_, file, line, ok := runtime.Caller(1)

	source := "unknown:0"
	if ok {
		source = filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	timestamp := time.Now().Format(time.RFC3339)

	_, err := fmt.Fprintf(os.Stdout, "[%s] [DEBUG] [%s] %s\n", timestamp, source, message)

	return err
}

func LogError(message string) error {
	_, err := fmt.Fprintln(os.Stderr, addTime(message))
	return err
}

func addTime(message string) string {
	return "[" + time.Now().Format(time.RFC3339) + "] " + message
}
