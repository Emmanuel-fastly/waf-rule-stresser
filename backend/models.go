// models.go defines the data structures used for test configuration and results.
package main

import "time"

// TestConfig represents the configuration sent from the frontend
// This struct defines what information we need to run a test
// The `json:"..."` tags tell Go how to convert between JSON and Go structs
type TestConfig struct {
	// TargetURL is the website/API we're testing
	// Example: "https://httpbin.org/get"
	TargetURL string `json:"target_url"`

	// TotalRequests is how many HTTP requests to send
	// Example: 100
	TotalRequests int `json:"total_requests"`

	// Duration is how long (in seconds) to spread the requests over
	// Example: 10 means send TotalRequests over 10 seconds
	Duration int `json:"duration"`

	// TrafficType determines what kind of requests to send
	// "normal" = regular requests
	// "attack" = requests with malicious payloads
	TrafficType string `json:"traffic_type"`

	// ErrorMode forces 404 responses by appending random paths
	// true = append /nonexistent-path-XXXXX to URL (tests 404 rate limits)
	// false = use URL as-is (tests 200 success rate limits)
	// Only used when TrafficType = "normal"
	ErrorMode bool `json:"error_mode"`

	// UserAgentType determines which user-agent pool to use
	// "legitimate" = real browser user-agents (Chrome, Firefox, etc.)
	// "scanner" = malicious scanner user-agents (sqlmap, Nikto, etc.)
	// Default: "legitimate" for normal traffic, "scanner" for attack traffic
	UserAgentType string `json:"user_agent_type"`

	// CustomUserAgent allows user to specify their own user-agent
	// This overrides UserAgentType if provided
	// Example: "MyCompany-Monitor/2.0"
	CustomUserAgent string `json:"custom_user_agent,omitempty"`

	// TestMode determines how requests are distributed over time
	// "baseline" = evenly distributed (1 req per second if 10 req in 10 sec)
	// "burst" = send faster at the beginning (stress test)
	TestMode string `json:"test_mode"`

	// HTTPMethod is the HTTP verb to use
	// Examples: "GET", "POST", "PUT", "DELETE", "PATCH"
	HTTPMethod string `json:"http_method"`

	// CustomHeaders are additional headers to include in each request
	// Example: {"Authorization": "Bearer token123", "X-API-Key": "secret"}
	CustomHeaders map[string]string `json:"custom_headers"`

	// RequestBody is the payload to send (for POST, PUT requests)
	// Example: `{"username": "test", "password": "test123"}`
	RequestBody string `json:"request_body"`
}

// RequestResult represents the outcome of a single HTTP request
// We create one of these for each request we send
type RequestResult struct {
	// ID is the request number (1, 2, 3, ...)
	ID int `json:"id"`

	// Status is the HTTP status code received
	// 200 = OK, 429 = Too Many Requests, 403 = Forbidden, etc.
	Status int `json:"status"`

	// StatusText is the human-readable status
	// Examples: "OK", "Too Many Requests", "Forbidden"
	StatusText string `json:"status_text"`

	// URL is the full URL we sent the request to
	// May include query parameters if attack mode is enabled
	URL string `json:"url"`

	// Method is the HTTP method used (GET, POST, etc.)
	Method string `json:"method"`

	// RequestHeaders are the headers we sent with this request
	RequestHeaders map[string]string `json:"request_headers"`

	// ResponseHeaders are the headers the server sent back
	ResponseHeaders map[string]string `json:"response_headers"`

	// RequestBody is what we sent (useful in attack mode to see the payload)
	RequestBody string `json:"request_body,omitempty"` // omitempty = don't include if empty

	// ResponseBody is what the server sent back (we'll truncate if too long)
	ResponseBody string `json:"response_body,omitempty"`

	// ResponseTime is how long the request took in milliseconds
	ResponseTime int `json:"response_time"`

	// Timestamp is exactly when this request was sent
	Timestamp time.Time `json:"timestamp"`

	// Error contains any error message if the request failed
	Error string `json:"error,omitempty"`

	// WasBlocked indicates if this request was rate-limited
	// true if status is 429 or 403
	WasBlocked bool `json:"was_blocked"`

	// AttackInfo describes what attack payload was used (if attack mode)
	// Example: "sql (query parameter)", "xss (request body)"
	AttackInfo string `json:"attack_info,omitempty"`
}

// TestResult represents the complete outcome of a test
// This is what we send back to the frontend when the test completes
type TestResult struct {
	// TestID is a unique identifier for this test
	// Format: "test_1234567890" (timestamp-based)
	TestID string `json:"test_id"`

	// TotalRequests is how many requests we actually sent
	TotalRequests int `json:"total_requests"`

	// SuccessCount is how many requests got 2xx status codes
	SuccessCount int `json:"success_count"`

	// ErrorCount is how many requests failed (4xx, 5xx, network errors)
	ErrorCount int `json:"error_count"`

	// BlockedCount is how many requests were rate-limited (429, 403)
	BlockedCount int `json:"blocked_count"`

	// AvgResponse is the average response time in milliseconds
	AvgResponse int `json:"avg_response"`

	// MinResponse is the fastest response time we got
	MinResponse int `json:"min_response"`

	// MaxResponse is the slowest response time we got
	MaxResponse int `json:"max_response"`

	// Percentiles for more detailed response time analysis
	P50Response int `json:"p50_response"` // Median (50th percentile)
	P95Response int `json:"p95_response"` // 95th percentile
	P99Response int `json:"p99_response"` // 99th percentile

	// Requests is the array of all individual request results
	Requests []RequestResult `json:"requests"`

	// RateLimitHit indicates if we detected rate limiting
	RateLimitHit bool `json:"rate_limit_hit"`

	// RateLimitAt is which request number triggered the rate limit
	// Example: 45 means the 45th request got blocked
	RateLimitAt int `json:"rate_limit_at"`

	// StartTime is when the test began
	StartTime time.Time `json:"start_time"`

	// EndTime is when the test completed
	EndTime time.Time `json:"end_time"`

	// Duration is the actual time taken in seconds
	Duration float64 `json:"duration"`

	// RequestsPerSec is the actual rate achieved
	RequestsPerSec float64 `json:"requests_per_sec"`
}

// ExportFormat is used when exporting test results to a file
// It includes both the configuration and the results
type ExportFormat struct {
	// Config is the test configuration that was used
	Config TestConfig `json:"config"`

	// Results is the outcome of the test
	Results TestResult `json:"results"`

	// ExportedAt is when this export was created
	ExportedAt time.Time `json:"exported_at"`
}

// ProgressUpdate represents a streaming update sent via SSE
type ProgressUpdate struct {
	// Type indicates the message type: "progress", "complete", "error", "cancelled"
	Type string `json:"type"`

	// TestID identifies which test this update belongs to
	TestID string `json:"test_id"`

	// Progress information (for "progress" type)
	Completed  int `json:"completed"`  // Number of requests completed
	Total      int `json:"total"`      // Total requests to send
	Percentage int `json:"percentage"` // Calculated percentage (0-100)

	// Latest batch of request results (sent every 1 second)
	NewRequests []RequestResult `json:"new_requests,omitempty"`

	// Incremental statistics (updated with each batch)
	CurrentStats *TestStats `json:"current_stats,omitempty"`

	// Final results (for "complete" type only)
	FinalResult *TestResult `json:"final_result,omitempty"`

	// Error information (for "error" type)
	Error string `json:"error,omitempty"`

	// Timestamp of this update
	Timestamp time.Time `json:"timestamp"`
}

// TestStats represents incremental statistics during test execution
type TestStats struct {
	SuccessCount int `json:"success_count"`
	ErrorCount   int `json:"error_count"`
	BlockedCount int `json:"blocked_count"`
	AvgResponse  int `json:"avg_response"`
}
