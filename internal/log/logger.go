package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
		source = fmt.Sprintf("%s:%d", filepath.Base(file), line)
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
	return fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC3339), message)
}
