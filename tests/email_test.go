package tests

import (
	db "github.com/Mans397/eLibrary/Database"
	"testing"
)

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{name: "Valid Gmail Email", email: "example@gmail.com", expected: true},
		{name: "Invalid Email without domain", email: "example", expected: false},
		{name: "Invalid Email with wrong domain", email: "example@yahoo.com", expected: false},
		{name: "Valid Gmail with extra text", email: "user@gmail.com.extra", expected: true},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := db.IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v; want %v", tt.email, result, tt.expected)
			} else {
				t.Logf("Test %d passed", i+1)
			}

		})
	}
}
