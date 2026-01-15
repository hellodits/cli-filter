package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"filter/internal/filter"
	"filter/internal/io"
	"filter/internal/parser"
)

const outputFile = "filtered_result.csv"

func main() {
	if err := run(); err != nil {
		fmt.Printf("unable to filter the data: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("successfully filter the data")
}

func run() error {
	// Parse CLI arguments
	dir := flag.String("d", "", "directory containing CSV files")
	startStr := flag.String("s", "", "start time (RFC3339, inclusive)")
	endStr := flag.String("e", "", "end time (RFC3339, exclusive)")
	flag.Parse()

	// Validate arguments
	if *dir == "" {
		return fmt.Errorf("missing required argument -d (directory)")
	}
	if *startStr == "" {
		return fmt.Errorf("missing required argument -s (start time)")
	}
	if *endStr == "" {
		return fmt.Errorf("missing required argument -e (end time)")
	}

	// Parse times
	start, err := parser.ParseTime(*startStr)
	if err != nil {
		return fmt.Errorf("invalid start time format: %v", err)
	}
	end, err := parser.ParseTime(*endStr)
	if err != nil {
		return fmt.Errorf("invalid end time format: %v", err)
	}

	// Validate directory exists
	info, err := os.Stat(*dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", *dir)
		}
		return fmt.Errorf("cannot access directory: %v", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", *dir)
	}

	// Get sorted CSV files
	files, err := io.GetSortedCSVFiles(*dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %v", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no CSV files found in directory: %s", *dir)
	}

	// Process files
	return processFiles(files, start, end)
}

func processFiles(files []string, start, end time.Time) error {
	writer, err := io.NewCSVWriter(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer writer.Close()

	for _, file := range files {
		shouldStop, err := processFile(file, start, end, writer)
		if err != nil {
			return fmt.Errorf("error processing %s: %v", file, err)
		}
		if shouldStop {
			break // Data is sequential, no need to process remaining files
		}
	}

	return nil
}

func processFile(path string, start, end time.Time, writer *io.CSVWriter) (bool, error) {
	reader, err := io.NewLineReader(path)
	if err != nil {
		return false, err
	}
	defer reader.Close()

	lineNum := 0
	for {
		line, ok := reader.ReadLine()
		if !ok {
			break
		}
		lineNum++

		if line == "" {
			continue // Skip empty lines
		}

		trx, err := parser.ParseLine(line)
		if err != nil {
			return false, fmt.Errorf("line %d: %v", lineNum, err)
		}

		result := filter.CheckRange(trx.TrxDate, start, end)
		switch result {
		case filter.Include:
			if err := writer.WriteLine(trx.Raw); err != nil {
				return false, fmt.Errorf("failed to write output: %v", err)
			}
		case filter.Stop:
			return true, nil // Stop processing all files
		case filter.Skip:
			// Continue to next line
		}
	}

	if err := reader.Err(); err != nil {
		return false, fmt.Errorf("read error: %v", err)
	}

	return false, nil
}
