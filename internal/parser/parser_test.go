package parser

import (
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid RFC3339 with timezone",
			input:   "2025-06-28T09:23:55+07:00",
			wantErr: false,
		},
		{
			name:    "valid RFC3339 UTC",
			input:   "2025-06-28T02:23:55Z",
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "2025-06-28 09:23:55",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "2025-13-28T09:23:55+07:00",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseTime(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTime() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseLine(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantTrxNo int
		wantErr   bool
	}{
		{
			name:      "valid line",
			input:     "84344,2025-06-28T09:23:55+07:00,transaction: 84344,1863012",
			wantTrxNo: 84344,
			wantErr:   false,
		},
		{
			name:      "empty line",
			input:     "",
			wantTrxNo: 0,
			wantErr:   true,
		},
		{
			name:      "missing columns",
			input:     "84344,2025-06-28T09:23:55+07:00",
			wantTrxNo: 0,
			wantErr:   true,
		},
		{
			name:      "invalid TrxNo",
			input:     "abc,2025-06-28T09:23:55+07:00,transaction: 84344,1863012",
			wantTrxNo: 0,
			wantErr:   true,
		},
		{
			name:      "invalid date",
			input:     "84344,invalid-date,transaction: 84344,1863012",
			wantTrxNo: 0,
			wantErr:   true,
		},
		{
			name:      "invalid amount",
			input:     "84344,2025-06-28T09:23:55+07:00,transaction: 84344,abc",
			wantTrxNo: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trx, err := ParseLine(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && trx.TrxNo != tt.wantTrxNo {
				t.Errorf("ParseLine() TrxNo = %v, want %v", trx.TrxNo, tt.wantTrxNo)
			}
		})
	}
}

func TestParseLinePreservesRaw(t *testing.T) {
	input := "84344,2025-06-28T09:23:55+07:00,transaction: 84344,1863012"
	trx, err := ParseLine(input)
	if err != nil {
		t.Fatalf("ParseLine() unexpected error: %v", err)
	}
	if trx.Raw != input {
		t.Errorf("ParseLine() Raw = %v, want %v", trx.Raw, input)
	}
}

func TestParseLineDate(t *testing.T) {
	input := "84344,2025-06-28T09:23:55+07:00,transaction: 84344,1863012"
	trx, err := ParseLine(input)
	if err != nil {
		t.Fatalf("ParseLine() unexpected error: %v", err)
	}

	expected := time.Date(2025, 6, 28, 9, 23, 55, 0, time.FixedZone("", 7*3600))
	if !trx.TrxDate.Equal(expected) {
		t.Errorf("ParseLine() TrxDate = %v, want %v", trx.TrxDate, expected)
	}
}
