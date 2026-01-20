package domain

import (
	"strings"
	"testing"
)

func TestValidateTodo(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid title",
			title:   "Buy milk",
			wantErr: false,
		},
		{
			name:    "empty title",
			title:   "",
			wantErr: true,
			errMsg:  "title cannot be empty",
		},
		{
			name:    "title too long",
			title:   strings.Repeat("a", 256),
			wantErr: true,
			errMsg:  "title too long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo := &Todo{Title: tt.title}
			err := ValidateTodo(todo)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTodo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("ValidateTodo() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
