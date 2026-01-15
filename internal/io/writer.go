package io

import (
	"bufio"
	"os"
)

// CSVWriter provides buffered writing to output file
type CSVWriter struct {
	file   *os.File
	writer *bufio.Writer
}

// NewCSVWriter creates a new writer, overwrites if file exists
func NewCSVWriter(path string) (*CSVWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &CSVWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// WriteLine writes a line to the output file
func (w *CSVWriter) WriteLine(line string) error {
	_, err := w.writer.WriteString(line + "\n")
	return err
}

// Close flushes and closes the writer
func (w *CSVWriter) Close() error {
	if err := w.writer.Flush(); err != nil {
		w.file.Close()
		return err
	}
	return w.file.Close()
}
