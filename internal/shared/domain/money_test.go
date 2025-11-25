package domain_test

import (
	"testing"

	"github.com/nguyenphuoc/super-salary-sacrifice/internal/shared/domain"
)

func TestMoney_NewMoney(t *testing.T) {
	tests := []struct {
		name    string
		dollars float64
		want    int64
		wantErr bool
	}{
		{"zero", 0, 0, false},
		{"positive", 10.50, 1050, false},
		{"large amount", 1000.99, 100099, false},
		{"negative", -10.00, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			money, err := domain.NewMoney(tt.dollars)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMoney() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && money.Cents() != tt.want {
				t.Errorf("Money.Cents() = %v, want %v", money.Cents(), tt.want)
			}
		})
	}
}

func TestMoney_Add(t *testing.T) {
	m1, _ := domain.NewMoney(10.50)
	m2, _ := domain.NewMoney(5.25)

	result := m1.Add(m2)

	if result.Dollars() != 15.75 {
		t.Errorf("Add() = %v, want 15.75", result.Dollars())
	}
}

func TestMoney_Subtract(t *testing.T) {
	m1, _ := domain.NewMoney(10.50)
	m2, _ := domain.NewMoney(5.25)

	result, err := m1.Subtract(m2)
	if err != nil {
		t.Errorf("Subtract() error = %v", err)
	}

	if result.Dollars() != 5.25 {
		t.Errorf("Subtract() = %v, want 5.25", result.Dollars())
	}

	// Test negative result
	_, err = m2.Subtract(m1)
	if err == nil {
		t.Error("Subtract() expected error for negative result")
	}
}

func TestMoney_Multiply(t *testing.T) {
	m, _ := domain.NewMoney(10.00)

	result := m.Multiply(2.5)

	if result.Dollars() != 25.00 {
		t.Errorf("Multiply() = %v, want 25.00", result.Dollars())
	}
}

func TestMoney_Comparisons(t *testing.T) {
	m1, _ := domain.NewMoney(10.00)
	m2, _ := domain.NewMoney(5.00)
	m3, _ := domain.NewMoney(10.00)

	if !m1.IsGreaterThan(m2) {
		t.Error("Expected m1 > m2")
	}

	if m2.IsGreaterThan(m1) {
		t.Error("Expected m2 < m1")
	}

	if !m1.Equals(m3) {
		t.Error("Expected m1 == m3")
	}

	if m1.Equals(m2) {
		t.Error("Expected m1 != m2")
	}
}

func TestMoney_Zero(t *testing.T) {
	zero := domain.Zero()

	if !zero.IsZero() {
		t.Error("Expected Zero() to return zero value")
	}

	if zero.Cents() != 0 {
		t.Errorf("Zero().Cents() = %v, want 0", zero.Cents())
	}
}
