package parser

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

// Transaction represents a single CSV record
type Transaction struct {
	TrxNo     int
	TrxDate   time.Time
	TrxDetail string
	Amount    int
	Raw       string // original line for output
}

var (
	ErrInvalidFormat = errors.New("invalid CSV format")
	ErrInvalidTrxNo  = errors.New("invalid TrxNo")
	ErrInvalidDate   = errors.New("invalid TrxDate")
	ErrInvalidAmount = errors.New("invalid Amount")
)

// ParseLine parses a single CSV line into Transaction
func ParseLine(line string) (Transaction, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return Transaction{}, ErrInvalidFormat
	}

	parts := strings.SplitN(line, ",", 4)
	if len(parts) != 4 {
		return Transaction{}, ErrInvalidFormat
	}

	trxNo, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return Transaction{}, ErrInvalidTrxNo
	}

	trxDate, err := time.Parse(time.RFC3339, strings.TrimSpace(parts[1]))
	if err != nil {
		return Transaction{}, ErrInvalidDate
	}

	amount, err := strconv.Atoi(strings.TrimSpace(parts[3]))
	if err != nil {
		return Transaction{}, ErrInvalidAmount
	}

	return Transaction{
		TrxNo:     trxNo,
		TrxDate:   trxDate,
		TrxDetail: parts[2],
		Amount:    amount,
		Raw:       line,
	}, nil
}

// ParseTime parses RFC3339 time string
func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
