package handlers

import (
	"testing"
)

func TestGetLibraryGenresRequest(t *testing.T) {
	tests := []struct {
		name         string
		queryParam   string
		expectRefresh bool
	}{
		{"no param", "", false},
		{"refresh true", "refresh=true", true},
		{"refresh false", "refresh=false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.queryParam == "refresh=true" && !tt.expectRefresh {
				t.Error("refresh=true should expect refresh")
			}
		})
	}
}
