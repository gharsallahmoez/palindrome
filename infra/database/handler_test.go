package database_test

import (
	"github.com/gharsallahmoez/palindrome/config"
	"github.com/gharsallahmoez/palindrome/infra/database"
	"testing"
)

// TestCreate tests Create function.
func TestCreate(t *testing.T) {
	// test table
	tt := []struct {
		name     string
		db       config.Database
		hasError bool
	}{
		{
			name: "valid database config",
			db: config.Database{
				Type: "in-memory",
			},
			hasError: false,
		},
		{
			name: "unsupported database config",
			db: config.Database{
				Type: "unsupported",
			},
			hasError: true,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := database.Create(tc.db)
			if err != nil && !tc.hasError {
				t.Errorf("expected success , got error: %v", err)
			}
			if err == nil && tc.hasError {
				t.Error("expected error, got nil")
			}
		})
	}
}
