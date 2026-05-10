package core_test

import (
	"testing"
	"time"

	"github.com/snivilised/jaywalk/src/agenor/core"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{name: "milliseconds", duration: 123 * time.Millisecond, expected: "123ms"},
		{name: "seconds", duration: 5 * time.Second, expected: "5s"},
		{name: "minutes and seconds", duration: 1*time.Minute + 2*time.Second, expected: "1:02"},
		{name: "full duration", duration: 2*time.Hour + 3*time.Minute + 4*time.Second, expected: "2h3m4s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := core.FormatDuration(tt.duration)
			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}
