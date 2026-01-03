package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// ExportResults saves test results to a file in the specified format
// format can be "json" or "csv"
// Returns the filename that was created
func ExportResults(config TestConfig, result TestResult, format string) (string, error) {
	// Create exports directory if it doesn't exist
	exportsDir := "../exports"
	err := os.MkdirAll(exportsDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create exports directory: %v", err)
	}

	// Generate filename with timestamp
	// Format: test_YYYYMMDD_HHMMSS.json or .csv
	timestamp := time.Now().Format("20060102_150405")
	var filename string

	switch format {
	case "json":
		filename = fmt.Sprintf("test_%s.json", timestamp)
		err = exportJSON(exportsDir, filename, config, result)
	case "csv":
		filename = fmt.Sprintf("test_%s.csv", timestamp)
		err = exportCSV(exportsDir, filename, config, result)
	default:
		return "", fmt.Errorf("unsupported format: %s (use 'json' or 'csv')", format)
	}

	if err != nil {
		return "", err
	}

	// Return the full path
	fullPath := filepath.Join(exportsDir, filename)
	return fullPath, nil
}

// exportJSON saves results in JSON format
// This preserves all data and structure
func exportJSON(dir, filename string, config TestConfig, result TestResult) error {
	// Create the export data structure
	exportData := ExportFormat{
		Config:     config,
		Results:    result,
		ExportedAt: time.Now(),
	}

	// Convert to JSON with pretty printing
	jsonData, err := json.MarshalIndent(exportData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	// Write to file
	fullPath := filepath.Join(dir, filename)
	err = os.WriteFile(fullPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON file: %v", err)
	}

	return nil
}

// exportCSV saves results in CSV format
// This is great for opening in Excel or Google Sheets
func exportCSV(dir, filename string, config TestConfig, result TestResult) error {
	fullPath := filepath.Join(dir, filename)

	// Create the CSV file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write test summary section
	writer.Write([]string{"TEST SUMMARY"})
	writer.Write([]string{"Test ID", result.TestID})
	writer.Write([]string{"Target URL", config.TargetURL})
	writer.Write([]string{"Total Requests", strconv.Itoa(result.TotalRequests)})
	writer.Write([]string{"Duration (seconds)", fmt.Sprintf("%.2f", result.Duration)})
	writer.Write([]string{"Traffic Type", config.TrafficType})
	writer.Write([]string{"Test Mode", config.TestMode})
	writer.Write([]string{"HTTP Method", config.HTTPMethod})
	writer.Write([]string{"Error Mode", strconv.FormatBool(config.ErrorMode)})
	writer.Write([]string{"User Agent Type", config.UserAgentType})
	if config.CustomUserAgent != "" {
		writer.Write([]string{"Custom User Agent", config.CustomUserAgent})
	}
	writer.Write([]string{}) // Empty row

	// Write statistics section
	writer.Write([]string{"STATISTICS"})
	writer.Write([]string{"Success Count", strconv.Itoa(result.SuccessCount)})
	writer.Write([]string{"Error Count", strconv.Itoa(result.ErrorCount)})
	writer.Write([]string{"Blocked Count", strconv.Itoa(result.BlockedCount)})
	writer.Write([]string{"Success Rate", fmt.Sprintf("%.1f%%", float64(result.SuccessCount)/float64(result.TotalRequests)*100)})
	writer.Write([]string{"Rate Limit Hit", strconv.FormatBool(result.RateLimitHit)})
	if result.RateLimitHit {
		writer.Write([]string{"Rate Limit At Request", strconv.Itoa(result.RateLimitAt)})
	}
	writer.Write([]string{"Requests Per Second", fmt.Sprintf("%.2f", result.RequestsPerSec)})
	writer.Write([]string{}) // Empty row

	// Write response time section
	writer.Write([]string{"RESPONSE TIMES (milliseconds)"})
	writer.Write([]string{"Average", strconv.Itoa(result.AvgResponse)})
	writer.Write([]string{"Minimum", strconv.Itoa(result.MinResponse)})
	writer.Write([]string{"Maximum", strconv.Itoa(result.MaxResponse)})
	writer.Write([]string{"p50 (Median)", strconv.Itoa(result.P50Response)})
	writer.Write([]string{"p95", strconv.Itoa(result.P95Response)})
	writer.Write([]string{"p99", strconv.Itoa(result.P99Response)})
	writer.Write([]string{}) // Empty row

	// Write individual requests section header
	writer.Write([]string{"INDIVIDUAL REQUESTS"})

	// CSV header row for requests
	headerRow := []string{
		"Request ID",
		"Timestamp",
		"Status Code",
		"Status Text",
		"Response Time (ms)",
		"Was Blocked",
		"URL",
		"Method",
		"User-Agent",
		"Error",
	}

	// Add attack info column if attack mode
	if config.TrafficType == "attack" {
		headerRow = append(headerRow, "Attack Info")
	}

	writer.Write(headerRow)

	// Write each request as a row
	for _, req := range result.Requests {
		row := []string{
			strconv.Itoa(req.ID),
			req.Timestamp.Format("2006-01-02 15:04:05"),
			strconv.Itoa(req.Status),
			req.StatusText,
			strconv.Itoa(req.ResponseTime),
			strconv.FormatBool(req.WasBlocked),
			req.URL,
			req.Method,
			req.RequestHeaders["User-Agent"],
			req.Error,
		}

		// Add attack info if present
		if config.TrafficType == "attack" {
			row = append(row, req.AttackInfo)
		}

		writer.Write(row)
	}

	return nil
}

// ListExports returns a list of all exported files
func ListExports() ([]string, error) {
	exportsDir := "../exports"

	// Check if directory exists
	if _, err := os.Stat(exportsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read directory contents
	files, err := os.ReadDir(exportsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read exports directory: %v", err)
	}

	// Collect filenames
	var exports []string
	for _, file := range files {
		if !file.IsDir() {
			exports = append(exports, file.Name())
		}
	}

	return exports, nil
}

// GetExportPath returns the full path to the exports directory
func GetExportPath() (string, error) {
	exportsDir := "../exports"
	absPath, err := filepath.Abs(exportsDir)
	if err != nil {
		return "", err
	}
	return absPath, nil
}
