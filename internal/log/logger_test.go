package log

import (
	"io"
	"os"
	"strings"
	"testing"
)

func captureOutput(fn func() error, captureStderr bool) (string, error) {
	if captureStderr {
		oldStderr := os.Stderr
		r, w, _ := os.Pipe()
		os.Stderr = w

		err := fn()

		w.Close()
		os.Stderr = oldStderr

		output, _ := io.ReadAll(r)
		return string(output), err
	} else {
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := fn()

		w.Close()
		os.Stdout = oldStdout

		output, _ := io.ReadAll(r)
		return string(output), err
	}
}

func TestLogInfo(t *testing.T) {
	message := "Test info message"

	output, err := captureOutput(func() error {
		return LogInfo(message)
	}, false)

	if err != nil {
		t.Errorf("LogInfo returned error: %v", err)
	}

	if !strings.Contains(output, message) {
		t.Errorf("Expected output to contain %q, got %q", message, output)
	}

	if !strings.Contains(output, "[") {
		t.Error("Expected output to contain timestamp")
	}
}

func TestLogError(t *testing.T) {
	message := "Test error message"

	output, err := captureOutput(func() error {
		return LogError(message)
	}, true)

	if err != nil {
		t.Errorf("LogError returned error: %v", err)
	}

	if !strings.Contains(output, message) {
		t.Errorf("Expected output to contain %q, got %q", message, output)
	}

	if !strings.Contains(output, "[") {
		t.Error("Expected output to contain timestamp")
	}
}

func TestLogDebug_Disabled(t *testing.T) {
	originalDebugMode := DebugMode
	defer func() { DebugMode = originalDebugMode }()

	DebugMode = false
	message := "Test debug message (should not appear)"

	output, err := captureOutput(func() error {
		return LogDebug(message)
	}, false)

	if err != nil {
		t.Errorf("LogDebug returned error: %v", err)
	}

	if strings.Contains(output, message) {
		t.Errorf("Expected no output when DebugMode is false, got %q", output)
	}
}

func TestLogDebug_Enabled(t *testing.T) {
	originalDebugMode := DebugMode
	defer func() { DebugMode = originalDebugMode }()

	DebugMode = true
	message := "Test debug message"

	output, err := captureOutput(func() error {
		return LogDebug(message)
	}, false)

	if err != nil {
		t.Errorf("LogDebug returned error: %v", err)
	}

	if !strings.Contains(output, message) {
		t.Errorf("Expected output to contain %q, got %q", message, output)
	}

	if !strings.Contains(output, "[DEBUG]") {
		t.Error("Expected output to contain [DEBUG]")
	}

	if !strings.Contains(output, "logger_test.go") {
		t.Error("Expected output to contain file source information")
	}
}

func TestLogDebug_ContainsLineNumber(t *testing.T) {
	originalDebugMode := DebugMode
	defer func() { DebugMode = originalDebugMode }()

	DebugMode = true

	output, _ := captureOutput(func() error {
		return LogDebug("test message")
	}, false)

	if !strings.Contains(output, "logger_test.go:") {
		t.Errorf("Expected output to contain file:line format, got %q", output)
	}
}

func TestAddTime_Format(t *testing.T) {
	message := "Test message"
	result := addTime(message)

	if !strings.Contains(result, "[") || !strings.Contains(result, "]") {
		t.Errorf("Expected timestamp in brackets, got %q", result)
	}

	if !strings.Contains(result, message) {
		t.Errorf("Expected message %q in result, got %q", message, result)
	}

	parts := strings.SplitN(result, "] ", 2)
	if len(parts) != 2 || !strings.HasSuffix(parts[1], message) {
		t.Errorf("Expected message to appear after timestamp, got %q", result)
	}
}

func TestAddTime_ConsistentFormat(t *testing.T) {
	messages := []string{"msg1", "msg2", "msg3"}

	for _, msg := range messages {
		result := addTime(msg)

		if strings.Count(result, "] ") != 1 {
			t.Errorf("Expected exactly one '] ' separator, got %q", result)
		}
	}
}

func TestDebugModeGlobalVariable(t *testing.T) {
	originalValue := DebugMode

	DebugMode = true
	if !DebugMode {
		t.Error("Failed to set DebugMode to true")
	}

	DebugMode = false
	if DebugMode {
		t.Error("Failed to set DebugMode to false")
	}

	DebugMode = originalValue
}

func TestLogInfo_EmptyMessage(t *testing.T) {
	output, err := captureOutput(func() error {
		return LogInfo("")
	}, false)

	if err != nil {
		t.Errorf("LogInfo returned error for empty message: %v", err)
	}

	if !strings.Contains(output, "[") {
		t.Error("Expected output to contain timestamp even for empty message")
	}
}

func TestLogError_EmptyMessage(t *testing.T) {
	output, err := captureOutput(func() error {
		return LogError("")
	}, true)

	if err != nil {
		t.Errorf("LogError returned error for empty message: %v", err)
	}

	if !strings.Contains(output, "[") {
		t.Error("Expected output to contain timestamp even for empty message")
	}
}

func TestLogDebug_EnabledDisabledBehavior(t *testing.T) {
	originalDebugMode := DebugMode
	defer func() { DebugMode = originalDebugMode }()

	message := "Debug test message"

	DebugMode = false
	output1, _ := captureOutput(func() error {
		return LogDebug(message)
	}, false)

	DebugMode = true
	output2, _ := captureOutput(func() error {
		return LogDebug(message)
	}, false)

	if strings.Contains(output1, message) {
		t.Error("Debug message should not appear when DebugMode is false")
	}

	if !strings.Contains(output2, message) {
		t.Error("Debug message should appear when DebugMode is true")
	}

	if strings.Contains(output1, "[DEBUG]") {
		t.Error("Should not contain [DEBUG] when DebugMode is false")
	}

	if !strings.Contains(output2, "[DEBUG]") {
		t.Error("Should contain [DEBUG] when DebugMode is true")
	}
}

func TestLogDebug_ReturnsNilWhenDisabled(t *testing.T) {
	originalDebugMode := DebugMode
	defer func() { DebugMode = originalDebugMode }()

	DebugMode = false
	err := LogDebug("test")

	if err != nil {
		t.Errorf("LogDebug should return nil when DebugMode is false, got %v", err)
	}
}
