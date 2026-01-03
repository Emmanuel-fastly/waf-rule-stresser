// tester.go contains the core logic for executing tests against the target URL.
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// RunTest is the main function that executes a test
// It takes a TestConfig and returns a TestResult
func RunTest(config TestConfig) TestResult {
	// Create a TestResult to store our findings
	result := TestResult{
		// Generate a unique test ID using the current timestamp
		TestID:    fmt.Sprintf("test_%d", time.Now().Unix()),
		StartTime: time.Now(),
		Requests:  make([]RequestResult, 0), // Initialize empty slice for requests
	}

	// Calculate the interval between requests based on test mode
	var intervalMs int

	if config.TestMode == "burst" {
		burstDuration := config.Duration / 2
		intervalMs = (burstDuration * 1000) / config.TotalRequests
	} else {
		intervalMs = (config.Duration * 1000) / config.TotalRequests
	}

	// Record exact start time for accurate statistics
	testStartTime := time.Now()

	// Send each request
	for i := 1; i <= config.TotalRequests; i++ {
		requestResult := sendRequest(i, config)

		// Add this result to our collection
		result.Requests = append(result.Requests, requestResult)

		// Check if this request was blocked (rate limited)
		if requestResult.WasBlocked && !result.RateLimitHit {
			result.RateLimitHit = true
			result.RateLimitAt = i
		}

		// Wait before sending the next request (except for the last one)
		if i < config.TotalRequests {
			time.Sleep(time.Duration(intervalMs) * time.Millisecond)
		}
	}

	// If burst mode, wait for the remaining time
	if config.TestMode == "burst" {
		elapsedSeconds := time.Since(testStartTime).Seconds()
		remainingSeconds := float64(config.Duration) - elapsedSeconds

		if remainingSeconds > 0 {
			time.Sleep(time.Duration(remainingSeconds * float64(time.Second)))
		}
	}

	// Calculate statistics
	result.EndTime = time.Now()
	result.TotalRequests = len(result.Requests)
	calculateStatistics(&result)

	return result
}

// sendRequest sends a single HTTP request and records the result
// This function handles all HTTP methods (GET, POST, PUT, DELETE)
func sendRequest(id int, config TestConfig) RequestResult {
	// Create the result struct for this request
	result := RequestResult{
		ID:        id,
		Timestamp: time.Now(),
		URL:       config.TargetURL,
		Method:    config.HTTPMethod,
	}

	// Default to GET if no method specified
	if result.Method == "" {
		result.Method = "GET"
	}

	// Prepare URL, body, and headers
	url := config.TargetURL
	body := config.RequestBody
	headers := make(map[string]string)

	// Copy custom headers from config
	for key, value := range config.CustomHeaders {
		headers[key] = value
	}

	// Handle normal traffic with error mode (404 testing)
	if config.TrafficType == "normal" && config.ErrorMode {
		url = Generate404URL(url)
		result.URL = url
	}

	// If attack mode is enabled, inject attack payloads
	if config.TrafficType == "attack" {
		url, body, headers, result.AttackInfo = BuildAttackRequest(config)
		result.URL = url
		result.RequestBody = body
	}

	// Create the HTTP request
	// We'll use bytes.NewBufferString to handle the request body
	var bodyReader io.Reader
	if body != "" {
		bodyReader = bytes.NewBufferString(body)
		if result.RequestBody == "" {
			result.RequestBody = body
		}
	}

	// Create the HTTP request object
	req, err := http.NewRequest(result.Method, url, bodyReader)
	if err != nil {
		// If we can't even create the request, record the error
		result.Error = fmt.Sprintf("Failed to create request: %v", err)
		result.Status = 0
		return result
	}

	// Add custom headers (including any attack headers)
	result.RequestHeaders = make(map[string]string)
	for key, value := range headers {
		req.Header.Set(key, value)
		result.RequestHeaders[key] = value
	}

	// Set Content-Type if we have a body
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
		result.RequestHeaders["Content-Type"] = "application/json"
	}

	// Select and set the appropriate User-Agent
	userAgent := SelectUserAgent(config)
	req.Header.Set("User-Agent", userAgent)
	result.RequestHeaders["User-Agent"] = userAgent

	// Create an HTTP client with a timeout
	// This prevents requests from hanging forever
	client := &http.Client{
		Timeout: 30 * time.Second, // 30 second timeout
	}

	// Record the start time so we can measure how long the request takes
	startTime := time.Now()

	// Send the request!
	resp, err := client.Do(req)

	// Calculate how long the request took
	elapsed := time.Since(startTime)
	result.ResponseTime = int(elapsed.Milliseconds())

	// Check if the request failed
	if err != nil {
		result.Error = fmt.Sprintf("Request failed: %v", err)
		result.Status = 0
		return result
	}

	// Make sure we close the response body when we're done
	// This is important to avoid memory leaks
	defer resp.Body.Close()

	// Record the status code
	result.Status = resp.StatusCode
	result.StatusText = http.StatusText(resp.StatusCode)

	// Check if this request was blocked (rate limited)
	// 429 = Too Many Requests, 406 = Not Acceptable, 403 = Forbidden (common WAF responses)
	if resp.StatusCode == 429 || resp.StatusCode == 406 || resp.StatusCode == 403 {
		result.WasBlocked = true
	}

	// Read the response headers
	result.ResponseHeaders = make(map[string]string)
	for key, values := range resp.Header {
		// Headers can have multiple values, we'll join them with commas
		result.ResponseHeaders[key] = strings.Join(values, ", ")
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = fmt.Sprintf("Failed to read response body: %v", err)
	} else {
		// Convert body to string
		bodyString := string(bodyBytes)

		// Truncate if too long (keep first 500 characters)
		if len(bodyString) > 500 {
			result.ResponseBody = bodyString[:500] + "... (truncated)"
		} else {
			result.ResponseBody = bodyString
		}
	}

	return result
}

// calculateStatistics computes the summary statistics for a test
// This modifies the TestResult in place (note the pointer receiver)
func calculateStatistics(result *TestResult) {
	// Initialize counters
	successCount := 0
	errorCount := 0
	blockedCount := 0
	totalResponseTime := 0
	minResponse := 999999 // Start with a very high number
	maxResponse := 0

	// Track response times for percentile calculations
	responseTimes := make([]int, 0, len(result.Requests))

	// Loop through all requests and count successes/errors
	for _, req := range result.Requests {
		// Count successes (2xx status codes)
		if req.Status >= 200 && req.Status < 300 {
			successCount++
		}

		// Count errors (4xx and 5xx, or network errors)
		if req.Status >= 400 || req.Error != "" {
			errorCount++
		}

		// Count blocked requests
		if req.WasBlocked {
			blockedCount++
		}

		// Sum up response times for average calculation
		totalResponseTime += req.ResponseTime
		responseTimes = append(responseTimes, req.ResponseTime)

		// Track min and max response times
		if req.ResponseTime < minResponse && req.ResponseTime > 0 {
			minResponse = req.ResponseTime
		}
		if req.ResponseTime > maxResponse {
			maxResponse = req.ResponseTime
		}
	}

	// Calculate average response time
	if result.TotalRequests > 0 {
		result.AvgResponse = totalResponseTime / result.TotalRequests
	}

	// Handle edge case: if no valid response times, set min to 0
	if minResponse == 999999 {
		minResponse = 0
	}

	// Set the statistics in the result
	result.SuccessCount = successCount
	result.ErrorCount = errorCount
	result.BlockedCount = blockedCount
	result.MinResponse = minResponse
	result.MaxResponse = maxResponse

	// Calculate percentiles for response time distribution
	p50, p95, p99 := calculateResponseTimePercentiles(responseTimes)
	result.P50Response = p50
	result.P95Response = p95
	result.P99Response = p99

	// Calculate actual duration and request rate
	result.Duration = result.EndTime.Sub(result.StartTime).Seconds()
	if result.Duration > 0 {
		result.RequestsPerSec = float64(result.TotalRequests) / result.Duration
	}

}

// calculateInterval extracts interval calculation logic
func calculateInterval(config TestConfig) int {
	if config.TestMode == "burst" {
		burstDuration := config.Duration / 2
		return (burstDuration * 1000) / config.TotalRequests
	}
	return (config.Duration * 1000) / config.TotalRequests
}

// calculateIncrementalStats computes stats from current requests
func calculateIncrementalStats(requests []RequestResult) *TestStats {
	stats := &TestStats{}
	totalResponseTime := 0

	for _, req := range requests {
		if req.Status >= 200 && req.Status < 300 {
			stats.SuccessCount++
		}
		if req.Status >= 400 || req.Error != "" {
			stats.ErrorCount++
		}
		if req.WasBlocked {
			stats.BlockedCount++
		}
		totalResponseTime += req.ResponseTime
	}

	if len(requests) > 0 {
		stats.AvgResponse = totalResponseTime / len(requests)
	}

	return stats
}

// buildFinalResult creates the complete TestResult from all requests
func buildFinalResult(testID string, requests []RequestResult, startTime time.Time) TestResult {
	result := TestResult{
		TestID:    testID,
		StartTime: startTime,
		EndTime:   time.Now(),
		Requests:  requests,
	}
	result.TotalRequests = len(requests)
	calculateStatistics(&result)
	return result
}

// sendRequestsConcurrently sends requests with proper timing and context cancellation
func sendRequestsConcurrently(ctx context.Context, config TestConfig, intervalMs int, resultsChan chan<- RequestResult) {
	defer close(resultsChan)

	for i := 1; i <= config.TotalRequests; i++ {
		// Check for cancellation before each request
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Send request
		result := sendRequest(i, config)

		// Try to send result (non-blocking check for cancellation)
		select {
		case <-ctx.Done():
			return
		case resultsChan <- result:
		}

		// Wait interval before next request (except for last one)
		if i < config.TotalRequests {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(intervalMs) * time.Millisecond):
			}
		}
	}
}

// RunTestStreaming executes a test with real-time progress updates
// Returns testID immediately; progress is sent via progressChan
func RunTestStreaming(ctx context.Context, config TestConfig, progressChan chan ProgressUpdate) string {
	testID := fmt.Sprintf("test_%d", time.Now().Unix())
	startTime := time.Now()

	// Send initial progress update
	progressChan <- ProgressUpdate{
		Type:       "progress",
		TestID:     testID,
		Completed:  0,
		Total:      config.TotalRequests,
		Percentage: 0,
		Timestamp:  time.Now(),
	}

	// Launch goroutine for test execution
	go func() {
		defer close(progressChan)

		// Channel to collect request results
		resultsChan := make(chan RequestResult, 100)

		// Batch collection variables
		var allRequests []RequestResult
		var batchRequests []RequestResult
		lastSendTime := time.Now()

		// Calculate interval
		intervalMs := calculateInterval(config)

		// Launch request sender goroutine
		go sendRequestsConcurrently(ctx, config, intervalMs, resultsChan)

		// Collect results and send batched updates
		completedCount := 0

		for result := range resultsChan {
			completedCount++
			allRequests = append(allRequests, result)
			batchRequests = append(batchRequests, result)

			// Check if context was cancelled
			select {
			case <-ctx.Done():
				progressChan <- ProgressUpdate{
					Type:      "cancelled",
					TestID:    testID,
					Completed: completedCount,
					Total:     config.TotalRequests,
					Timestamp: time.Now(),
				}
				return
			default:
			}

			// Send batch every 1 second OR when test completes
			timeSinceLastSend := time.Since(lastSendTime)
			isComplete := completedCount >= config.TotalRequests

			if timeSinceLastSend >= 1*time.Second || isComplete {
				// Calculate incremental stats
				stats := calculateIncrementalStats(allRequests)

				// Calculate percentage
				percentage := (completedCount * 100) / config.TotalRequests

				// Send progress update with batch
				progressChan <- ProgressUpdate{
					Type:         "progress",
					TestID:       testID,
					Completed:    completedCount,
					Total:        config.TotalRequests,
					Percentage:   percentage,
					NewRequests:  batchRequests,
					CurrentStats: stats,
					Timestamp:    time.Now(),
				}

				// Reset batch
				batchRequests = []RequestResult{}
				lastSendTime = time.Now()
			}

			if isComplete {
				break
			}
		}

		// Calculate final statistics
		finalResult := buildFinalResult(testID, allRequests, startTime)

		// Send completion message
		progressChan <- ProgressUpdate{
			Type:        "complete",
			TestID:      testID,
			Completed:   config.TotalRequests,
			Total:       config.TotalRequests,
			Percentage:  100,
			FinalResult: &finalResult,
			Timestamp:   time.Now(),
		}
	}()

	return testID
}
