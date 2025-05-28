package utils

import (
	"testing"

	"github.com/spf13/viper"
)

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
	}{
		{"default info", false},
		{"verbose debug", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Set("verbose", tt.verbose)
			// Should not panic
			InitLogger()
		})
	}
}
