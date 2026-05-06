package util

import "testing"

func TestValidateText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "Buy milk", wantErr: false},
		{name: "trimmed empty", input: "   ", wantErr: true},
		{name: "too long", input: string(make([]byte, 251)), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateText(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ValidateText() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseDueDate(t *testing.T) {
	if _, err := ParseDueDate("2026-05-10"); err != nil {
		t.Fatalf("expected valid date, got error: %v", err)
	}
	if _, err := ParseDueDate("10-05-2026"); err == nil {
		t.Fatal("expected invalid format error")
	}
}
