package validation_test

import (
	"testing"

	"github.com/kripsy/GophKeeper/internal/client/infrastrucrure/ui/validation"
)

func Test_validateUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "ok",
			input:   "password",
			wantErr: false,
		},
		{
			name:    "invalid pass less then 6",
			input:   "pass",
			wantErr: true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validation.ValidateCardNumber(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateUsername() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "ok",
			input:   "username",
			wantErr: false,
		},
		{
			name:    "invalid username less then 6",
			input:   "login",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validation.ValidatePassword(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateCVV(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "ok",
			input:   "123",
			wantErr: false,
		},
		{
			name:    "invalid cvv not num",
			input:   "cvv",
			wantErr: true,
		},
		{
			name:    "invalid cvv less then 3",
			input:   "12",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validation.ValidateCVV(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateCVV() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateCardNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "ok",
			input:   "123123123123",
			wantErr: false,
		},
		{
			name:    "invalid card num not num",
			input:   "card",
			wantErr: true,
		},
		{
			name:    "invalid card bigger then 16",
			input:   "12345678901234567",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validation.ValidateCardNumber(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("validateCardNumber() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
