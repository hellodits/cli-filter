package io

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractPrefix(t *testing.T) {
	tests := []struct {
		filename string
		expected int
	}{
		{"1_report.csv", 1},
		{"10_report.csv", 10},
		{"100_data.csv", 100},
		{"report.csv", 0},
		{"_report.csv", 0},
		{"abc_report.csv", 0},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := extractPrefix(tt.filename)
			if result != tt.expected {
				t.Errorf("extractPrefix(%s) = %d, want %d", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestGetSortedCSVFiles(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "filter_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files in random order
	files := []string{"3_report.csv", "1_report.csv", "10_report.csv", "2_report.csv"}
	for _, f := range files {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Get sorted files
	result, err := GetSortedCSVFiles(tmpDir)
	if err != nil {
		t.Fatalf("GetSortedCSVFiles() error: %v", err)
	}

	// Verify order
	expected := []string{"1_report.csv", "2_report.csv", "3_report.csv", "10_report.csv"}
	if len(result) != len(expected) {
		t.Fatalf("GetSortedCSVFiles() returned %d files, want %d", len(result), len(expected))
	}

	for i, exp := range expected {
		if filepath.Base(result[i]) != exp {
			t.Errorf("GetSortedCSVFiles()[%d] = %s, want %s", i, filepath.Base(result[i]), exp)
		}
	}
}

func TestLineReader(t *testing.T) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "filter_test_*.csv")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := "line1\nline2\nline3"
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write test content: %v", err)
	}
	tmpFile.Close()

	// Test reading
	reader, err := NewLineReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewLineReader() error: %v", err)
	}
	defer reader.Close()

	expected := []string{"line1", "line2", "line3"}
	for i, exp := range expected {
		line, ok := reader.ReadLine()
		if !ok {
			t.Fatalf("ReadLine() returned false at line %d", i+1)
		}
		if line != exp {
			t.Errorf("ReadLine() = %s, want %s", line, exp)
		}
	}

	// Should return false at EOF
	_, ok := reader.ReadLine()
	if ok {
		t.Error("ReadLine() should return false at EOF")
	}
}
