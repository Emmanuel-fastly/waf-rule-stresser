// handlers.go contains the HTTP handlers for the API endpoints.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// TestSession manages an active test's execution and cancellation
type TestSession struct {
	TestID       string
	Config       TestConfig
	CancelFunc   context.CancelFunc
	ProgressChan chan ProgressUpdate
	StartTime    time.Time
}

var (
	// activeSessions stores all currently running test sessions
	activeSessions = make(map[string]*TestSession)

	// sessionsMutex protects concurrent access to activeSessions
	sessionsMutex sync.RWMutex
)

// registerSession adds a session to the active sessions map
func registerSession(session *TestSession) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	activeSessions[session.TestID] = session
}

// unregisterSession removes a session from the active sessions map
func unregisterSession(testID string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	delete(activeSessions, testID)
}

// getSession retrieves a session from the active sessions map
func getSession(testID string) (*TestSession, bool) {
	sessionsMutex.RLock()
	defer sessionsMutex.RUnlock()
	session, exists := activeSessions[testID]
	return session, exists
}

// handleStartTest handles POST requests to /api/test/start
// This is the main endpoint that receives test configuration and runs the test
func handleStartTest(w http.ResponseWriter, r *http.Request) {
	// 1. VERIFY THE REQUEST METHOD
	// We only accept POST requests for starting tests
	if r.Method != http.MethodPost {
		// If it's not POST, send an error response
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	// Enable CORS
	enableCORS(w, r)

	// Handle preflight OPTIONS request
	// Browsers send this before the actual POST request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 3. PARSE THE REQUEST BODY
	// The request body contains JSON with the test configuration
	var config TestConfig

	// Create a JSON decoder that reads from the request body
	decoder := json.NewDecoder(r.Body)

	// Decode the JSON into our TestConfig struct
	err := decoder.Decode(&config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// 4. VALIDATE THE CONFIGURATION
	// Make sure the config has all required fields
	if config.TargetURL == "" {
		http.Error(w, "target_url is required", http.StatusBadRequest)
		return
	}
	if config.TotalRequests <= 0 {
		http.Error(w, "total_requests must be greater than 0", http.StatusBadRequest)
		return
	}
	if config.TotalRequests > 10000 {
		http.Error(w, "total_requests cannot exceed 10000 (safety limit)", http.StatusBadRequest)
		return
	}
	if config.Duration <= 0 {
		http.Error(w, "duration must be greater than 0", http.StatusBadRequest)
		return
	}
	if config.Duration > 3600 {
		http.Error(w, "duration cannot exceed 3600 seconds (1 hour)", http.StatusBadRequest)
		return
	}

	// Validate traffic type
	if config.TrafficType != "" && config.TrafficType != "normal" && config.TrafficType != "attack" {
		http.Error(w, "traffic_type must be 'normal' or 'attack'", http.StatusBadRequest)
		return
	}

	// Validate test mode
	if config.TestMode != "" && config.TestMode != "baseline" && config.TestMode != "burst" {
		http.Error(w, "test_mode must be 'baseline' or 'burst'", http.StatusBadRequest)
		return
	}

	// Validate user agent type
	if config.UserAgentType != "" && config.UserAgentType != "legitimate" && config.UserAgentType != "scanner" {
		http.Error(w, "user_agent_type must be 'legitimate' or 'scanner'", http.StatusBadRequest)
		return
	}

	// Set defaults if not provided
	if config.TrafficType == "" {
		config.TrafficType = "normal"
	}
	if config.TestMode == "" {
		config.TestMode = "baseline"
	}
	if config.HTTPMethod == "" {
		config.HTTPMethod = "GET"
	}

	// Run the test
	result := RunTest(config)

	// 6. SEND THE RESPONSE
	// Set the response header to indicate we're sending JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the result as JSON and send it
	encoder := json.NewEncoder(w)
	// PrettyPrint the JSON (makes it readable for debugging)
	encoder.SetIndent("", "  ")

	err = encoder.Encode(result)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// enableCORS sets the necessary headers for Cross-Origin Resource Sharing
func enableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

// handleExport handles POST requests to /api/test/export
// This exports test results to JSON or CSV format
func handleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	enableCORS(w, r)

	// Handle preflight
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse the request body
	// Expected: { config: {...}, results: {...}, format: "json" or "csv" }
	var requestData struct {
		Config  TestConfig `json:"config"`
		Results TestResult `json:"results"`
		Format  string     `json:"format"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	if requestData.Format != "json" && requestData.Format != "csv" {
		http.Error(w, "format must be 'json' or 'csv'", http.StatusBadRequest)
		return
	}

	filename, err := ExportResults(requestData.Config, requestData.Results, requestData.Format)
	if err != nil {
		http.Error(w, fmt.Sprintf("Export failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Send success response with filename
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":   "success",
		"filename": filename,
		"message":  fmt.Sprintf("Results exported to %s", filename),
	}

	json.NewEncoder(w).Encode(response)
}

// handleListExports handles GET requests to /api/exports/list
// Returns a list of all exported files
func handleListExports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed. Use GET.", http.StatusMethodNotAllowed)
		return
	}

	enableCORS(w, r)

	exports, err := ListExports()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list exports: %v", err), http.StatusInternalServerError)
		return
	}

	// Get exports directory path
	exportPath, _ := GetExportPath()

	// Send response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"exports":     exports,
		"count":       len(exports),
		"export_path": exportPath,
	}

	json.NewEncoder(w).Encode(response)
}

// handleStartTestStream handles streaming test execution via SSE
func handleStartTestStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	enableCORS(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var config TestConfig
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&config)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// 3. Validate configuration (same as handleStartTest)
	if config.TargetURL == "" {
		http.Error(w, "target_url is required", http.StatusBadRequest)
		return
	}
	if config.TotalRequests <= 0 {
		http.Error(w, "total_requests must be greater than 0", http.StatusBadRequest)
		return
	}
	if config.TotalRequests > 10000 {
		http.Error(w, "total_requests cannot exceed 10000 (safety limit)", http.StatusBadRequest)
		return
	}
	if config.Duration <= 0 {
		http.Error(w, "duration must be greater than 0", http.StatusBadRequest)
		return
	}
	if config.Duration > 3600 {
		http.Error(w, "duration cannot exceed 3600 seconds (1 hour)", http.StatusBadRequest)
		return
	}

	// Validate traffic type
	if config.TrafficType != "" && config.TrafficType != "normal" && config.TrafficType != "attack" {
		http.Error(w, "traffic_type must be 'normal' or 'attack'", http.StatusBadRequest)
		return
	}

	// Validate test mode
	if config.TestMode != "" && config.TestMode != "baseline" && config.TestMode != "burst" {
		http.Error(w, "test_mode must be 'baseline' or 'burst'", http.StatusBadRequest)
		return
	}

	// Validate user agent type
	if config.UserAgentType != "" && config.UserAgentType != "legitimate" && config.UserAgentType != "scanner" {
		http.Error(w, "user_agent_type must be 'legitimate' or 'scanner'", http.StatusBadRequest)
		return
	}

	// Set defaults if not provided
	if config.TrafficType == "" {
		config.TrafficType = "normal"
	}
	if config.TestMode == "" {
		config.TestMode = "baseline"
	}
	if config.HTTPMethod == "" {
		config.HTTPMethod = "GET"
	}

	// Set up SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// 5. Create flusher for real-time streaming
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	// 6. Create context with cancellation
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// 7. Create progress channel
	progressChan := make(chan ProgressUpdate, 10)

	// 8. Start test in goroutine
	testID := RunTestStreaming(ctx, config, progressChan)

	// 9. Register session for cancellation support
	session := &TestSession{
		TestID:       testID,
		Config:       config,
		CancelFunc:   cancel,
		ProgressChan: progressChan,
		StartTime:    time.Now(),
	}
	registerSession(session)
	defer unregisterSession(testID)

	// Stream updates to client
	for update := range progressChan {
		jsonData, err := json.Marshal(update)
		if err != nil {
			continue
		}

		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	}
}

// handleStopTest handles POST requests to cancel a running test
func handleStopTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed. Use POST.", http.StatusMethodNotAllowed)
		return
	}

	enableCORS(w, r)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Parse request body to get testID
	var requestData struct {
		TestID string `json:"test_id"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&requestData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Get session and cancel
	session, exists := getSession(requestData.TestID)
	if !exists {
		http.Error(w, "Test not found or already completed", http.StatusNotFound)
		return
	}

	// Cancel the context (this stops ongoing requests)
	session.CancelFunc()

	// Send success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "cancelled",
		"test_id": requestData.TestID,
	})
}
