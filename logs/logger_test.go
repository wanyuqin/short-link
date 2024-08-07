package logs

import (
	"errors"
	"testing"

	"go.uber.org/zap"
)

func TestError(t *testing.T) {
	Error(errors.New("xxx"), "aaa %s %d", zap.Any("a", "a"))
}

func TestInitializeLogger(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitializeLogger()
		})
	}
}
