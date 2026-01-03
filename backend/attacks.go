// attacks.go contains logic for loading and managing attack payloads
package main

import (
	"bufio"
	"embed"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

//go:embed payloads/*.txt
var payloadsFS embed.FS

// AttackPayloads stores all loaded attack payloads in memory
// This is a global variable that gets populated when the server starts
var AttackPayloads = struct {
	SQLInjection         []string
	XSS                  []string
	PathTraversal        []string
	CommandInjection     []string
	ScannerUserAgents    []string
	LegitimateUserAgents []string
}{
	// Initialize with empty slices
	SQLInjection:         []string{},
	XSS:                  []string{},
	PathTraversal:        []string{},
	CommandInjection:     []string{},
	ScannerUserAgents:    []string{},
	LegitimateUserAgents: []string{},
}

// LoadPayloads reads all attack payload files and stores them in memory
// This should be called once when the server starts
// Returns an error if any file fails to load
func LoadPayloads() error {
	// Define the payload files we need to load
	payloadFiles := map[string]*[]string{
		"SQL Injection":          &AttackPayloads.SQLInjection,
		"XSS":                    &AttackPayloads.XSS,
		"Path Traversal":         &AttackPayloads.PathTraversal,
		"Command Injection":      &AttackPayloads.CommandInjection,
		"Scanner User-Agents":    &AttackPayloads.ScannerUserAgents,
		"Legitimate User-Agents": &AttackPayloads.LegitimateUserAgents,
	}

	filenames := map[string]string{
		"SQL Injection":          "payloads/sql-injection.txt",
		"XSS":                    "payloads/xss.txt",
		"Path Traversal":         "payloads/path-traversal.txt",
		"Command Injection":      "payloads/command-injection.txt",
		"Scanner User-Agents":    "payloads/scanner-user-agents.txt",
		"Legitimate User-Agents": "payloads/legitimate-user-agents.txt",
	}

	// Load each payload file
	for name, target := range payloadFiles {
		filename := filenames[name]

		// Load the file
		payloads, err := loadPayloadFile(filename)
		if err != nil {
			// If file doesn't exist or can't be read, skip silently
			continue
		}

		// Store the loaded payloads
		*target = payloads
	}

	// Check if we loaded at least some payloads
	totalPayloads := len(AttackPayloads.SQLInjection) +
		len(AttackPayloads.XSS) +
		len(AttackPayloads.PathTraversal) +
		len(AttackPayloads.CommandInjection) +
		len(AttackPayloads.ScannerUserAgents) +
		len(AttackPayloads.LegitimateUserAgents)

	if totalPayloads == 0 {
		return fmt.Errorf("no payloads loaded")
	}

	return nil
}

// loadPayloadFile reads a single payload file and returns the payloads as a slice
// Each line in the file becomes one payload (blank lines and comments are skipped)
func loadPayloadFile(filename string) ([]string, error) {
	// Read from embedded filesystem
	data, err := payloadsFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Create a slice to store the payloads
	payloads := []string{}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(strings.NewReader(string(data)))

	// Read each line
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines
		if line == "" {
			continue
		}

		// Skip comment lines (starting with #)
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Add this payload to our collection
		payloads = append(payloads, line)
	}

	return payloads, nil
}

// GetRandomPayload returns a random payload from the specified type
// attackType can be: "sql", "xss", "traversal", "command"
// Returns empty string if no payloads available for that type
func GetRandomPayload(attackType string) string {
	var payloads []string

	// Select the appropriate payload collection
	switch attackType {
	case "sql":
		payloads = AttackPayloads.SQLInjection
	case "xss":
		payloads = AttackPayloads.XSS
	case "traversal":
		payloads = AttackPayloads.PathTraversal
	case "command":
		payloads = AttackPayloads.CommandInjection
	default:
		return ""
	}

	// If no payloads loaded for this type, return empty string
	if len(payloads) == 0 {
		return ""
	}

	// Return a random payload from the collection
	randomIndex := rand.Intn(len(payloads))
	return payloads[randomIndex]
}

// GetRandomUserAgent returns a random malicious user-agent string
// These are commonly blocked by WAFs
func GetRandomUserAgent() string {
	if len(AttackPayloads.ScannerUserAgents) == 0 {
		// Fallback user agents if file not loaded
		fallbackAgents := []string{
			"sqlmap/1.0",
			"Nikto/2.1.6",
			"nmap-scripting-engine",
		}
		randomIndex := rand.Intn(len(fallbackAgents))
		return fallbackAgents[randomIndex]
	}

	// Return a random user agent from the loaded collection
	randomIndex := rand.Intn(len(AttackPayloads.ScannerUserAgents))
	return AttackPayloads.ScannerUserAgents[randomIndex]
}

// GetLegitimateUserAgent returns a random legitimate browser user-agent
// These are normal browser signatures that should not be blocked
func GetLegitimateUserAgent() string {
	if len(AttackPayloads.LegitimateUserAgents) == 0 {
		// Fallback to latest Chrome if file not loaded
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36"
	}

	// Return a random legitimate user agent from the loaded collection
	randomIndex := rand.Intn(len(AttackPayloads.LegitimateUserAgents))
	return AttackPayloads.LegitimateUserAgents[randomIndex]
}

// SelectUserAgent returns the appropriate user-agent based on config
// Priority: CustomUserAgent > UserAgentType selection > defaults
func SelectUserAgent(config TestConfig) string {
	// Priority 1: Custom user-agent (overrides everything)
	if config.CustomUserAgent != "" {
		return config.CustomUserAgent
	}

	// Priority 2: Explicit user-agent type selection
	if config.UserAgentType == "scanner" {
		return GetRandomUserAgent()
	}

	if config.UserAgentType == "legitimate" {
		return GetLegitimateUserAgent()
	}

	// Priority 3: Defaults based on traffic type
	if config.TrafficType == "attack" {
		// Attack traffic defaults to scanner UAs
		return GetRandomUserAgent()
	}

	// Normal traffic defaults to legitimate UAs
	return GetLegitimateUserAgent()
}

// GetRandomAttackType returns a random attack type for variety
// This ensures we test different attack vectors
func GetRandomAttackType() string {
	attackTypes := []string{"sql", "xss", "traversal", "command"}
	randomIndex := rand.Intn(len(attackTypes))
	return attackTypes[randomIndex]
}

// BuildAttackRequest modifies a request to include attack payloads
// This is called from tester.go when traffic_type is "attack"
// It randomly selects an attack type and injects the payload
func BuildAttackRequest(config TestConfig) (url string, body string, headers map[string]string, attackInfo string) {
	// Start with the original URL and body
	url = config.TargetURL
	body = config.RequestBody

	// Copy custom headers
	headers = make(map[string]string)
	for key, value := range config.CustomHeaders {
		headers[key] = value
	}

	// Choose a random attack type
	attackType := GetRandomAttackType()
	payload := GetRandomPayload(attackType)

	// If we couldn't get a payload, just return original request
	if payload == "" {
		attackInfo = "No payloads available"
		return
	}

	// Inject the payload based on the HTTP method
	if config.HTTPMethod == "GET" || config.HTTPMethod == "" {
		// For GET requests, add payload as query parameter
		separator := "?"
		if strings.Contains(url, "?") {
			separator = "&"
		}
		url = fmt.Sprintf("%s%sid=%s&search=%s", url, separator, payload, payload)
		attackInfo = fmt.Sprintf("%s (query parameter)", attackType)

	} else {
		// For POST/PUT/DELETE, inject in body
		// If body is empty, create a JSON body with the payload
		if body == "" {
			body = fmt.Sprintf(`{"id": "%s", "search": "%s"}`, payload, payload)
		} else {
			// If body exists, try to inject payload into it
			// For simplicity, we'll append to the body
			body = strings.TrimSuffix(body, "}")
			body = fmt.Sprintf(`%s, "attack": "%s"}`, body, payload)
		}
		attackInfo = fmt.Sprintf("%s (request body)", attackType)
	}

	// Always use a malicious user-agent for attack requests
	headers["User-Agent"] = GetRandomUserAgent()

	return url, body, headers, attackInfo
}

// Generate404URL modifies a URL to force a 404 response
// This is used to test rate limits on error responses
func Generate404URL(baseURL string) string {
	// Generate a random path that definitely won't exist
	randomPath := fmt.Sprintf("/nonexistent-path-%d-%d", time.Now().Unix(), rand.Intn(999999))

	// Remove trailing slash from base URL if present
	baseURL = strings.TrimSuffix(baseURL, "/")

	// Append the nonexistent path
	return baseURL + randomPath
}
