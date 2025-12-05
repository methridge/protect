package logger

import (
	"testing"
)

func TestNew(t *testing.T) {
	globalLogger = nil

	logger := New()
	if logger == nil {
		t.Error("Expected logger to be initialized")
	}

	logger2 := New()
	if logger != logger2 {
		t.Error("Expected New() to return the same logger instance")
	}
}

func TestGet(t *testing.T) {
	globalLogger = nil

	logger := Get()
	if logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestSetLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{"debug level", "debug"},
		{"info level", "info"},
		{"warn level", "warn"},
		{"error level", "error"},
		{"default level", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetLevel(tt.level)
			if err != nil {
				t.Errorf("SetLevel(%s) error = %v", tt.level, err)
			}

			logger := Get()
			if logger == nil {
				t.Error("Expected logger to be initialized after SetLevel")
			}
		})
	}
}
