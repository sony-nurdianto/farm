package uni_test

import (
	"testing"

	"github.com/sony-nurdianto/farm/auth/internal/validator"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email string
		want  bool
	}{
		{"test@example.com", true},
		{"user.name+tag+sorting@example.com", true},
		{"invalid-email", false},
		{"another.invalid@com", false},
	}

	for _, tt := range tests {
		got := validator.ValidateEmail(tt.email)
		if got != tt.want {
			t.Errorf("ValidateEmail(%q) = %v; want %v", tt.email, got, tt.want)
		}
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		phone string
		want  bool
	}{
		{"081234567890", true},
		{"+6281234567890", true},
		{"12345", false},
		{"phone123456", false},
	}

	for _, tt := range tests {
		got := validator.ValidatePhone(tt.phone)
		if got != tt.want {
			t.Errorf("ValidatePhone(%q) = %v; want %v", tt.phone, got, tt.want)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		password string
		want     bool
	}{
		{"Abcdef1!", true},
		{"abcdefg1!", false}, // no uppercase
		{"ABCDEFGH!", false}, // no number
		{"Abcdefgh1", false}, // no special char
		{"Ab1!", false},      // too short
	}

	for _, tt := range tests {
		got := validator.ValidatePassword(tt.password)
		if got != tt.want {
			t.Errorf("ValidatePassword(%q) = %v; want %v", tt.password, got, tt.want)
		}
	}
}
