package domain_test

import (
	"testing"

	"github.com/nguyenphuoc/super-salary-sacrifice/internal/shared/domain"
)

func TestEmail_Valid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"valid email", "test@example.com", "test@example.com"},
		{"uppercase", "TEST@EXAMPLE.COM", "test@example.com"},
		{"with spaces", "  test@example.com  ", "test@example.com"},
		{"with plus", "test+tag@example.com", "test+tag@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := domain.NewEmail(tt.input)
			if err != nil {
				t.Errorf("NewEmail() error = %v, want nil", err)
				return
			}
			if email.String() != tt.want {
				t.Errorf("Email.String() = %v, want %v", email.String(), tt.want)
			}
		})
	}
}

func TestEmail_Invalid(t *testing.T) {
	tests := []string{
		"",
		"invalid",
		"@example.com",
		"test@",
		"test @example.com",
		"test@example",
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := domain.NewEmail(input)
			if err == nil {
				t.Errorf("NewEmail(%q) expected error, got nil", input)
			}
		})
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, _ := domain.NewEmail("test@example.com")
	email2, _ := domain.NewEmail("TEST@EXAMPLE.COM")
	email3, _ := domain.NewEmail("other@example.com")

	if !email1.Equals(email2) {
		t.Error("Expected emails to be equal (case-insensitive)")
	}

	if email1.Equals(email3) {
		t.Error("Expected emails to be different")
	}
}
