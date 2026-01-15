package io

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// GetSortedCSVFiles returns CSV files sorted by numeric prefix
func GetSortedCSVFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".csv") {
			files = append(files, filepath.Join(dir, name))
		}
	}

	// Sort by numeric prefix
	sort.Slice(files, func(i, j int) bool {
		return extractPrefix(filepath.Base(files[i])) < extractPrefix(filepath.Base(files[j]))
	})

	return files, nil
}

// extractPrefix extracts numeric prefix from filename (e.g., "1_report.csv" -> 1)
func extractPrefix(filename string) int {
	idx := strings.Index(filename, "_")
	if idx == -1 {
		return 0
	}
	n, err := strconv.Atoi(filename[:idx])
	if err != nil {
		return 0
	}
	return n
}

// LineReader provides streaming line-by-line reading
type LineReader struct {
	file    *os.File
	scanner *bufio.Scanner
}

// NewLineReader creates a new streaming reader for a file
func NewLineReader(path string) (*LineReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &LineReader{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

// ReadLine reads the next line, returns empty string and false at EOF
func (r *LineReader) ReadLine() (string, bool) {
	if r.scanner.Scan() {
		return r.scanner.Text(), true
	}
	return "", false
}

// Err returns any error from scanning
func (r *LineReader) Err() error {
	return r.scanner.Err()
}

// Close closes the underlying file
func (r *LineReader) Close() error {
	return r.file.Close()
}
