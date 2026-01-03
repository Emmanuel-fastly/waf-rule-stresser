package main

import "sort"

// calculatePercentile calculates a specific percentile from a sorted slice of integers
// percentile should be between 0 and 100 (e.g., 50 for median, 95 for p95)
// This uses the "nearest rank" method
func calculatePercentile(sortedValues []int, percentile float64) int {
	if len(sortedValues) == 0 {
		return 0
	}

	// Calculate the index for this percentile
	// Formula: index = (percentile / 100) * (n - 1)
	// Example: For p95 with 100 values: (95/100) * 99 = 94.05 → index 94
	index := (percentile / 100.0) * float64(len(sortedValues)-1)

	// Round to nearest integer
	// We use int() which truncates, but add 0.5 first for proper rounding
	roundedIndex := int(index + 0.5)

	// Make sure we don't go out of bounds
	if roundedIndex >= len(sortedValues) {
		roundedIndex = len(sortedValues) - 1
	}

	return sortedValues[roundedIndex]
}

// sortInts sorts a slice of integers in ascending order
// This is a helper function that makes a copy to avoid modifying the original
func sortInts(values []int) []int {
	// Make a copy so we don't modify the original slice
	sorted := make([]int, len(values))
	copy(sorted, values)

	// Sort in ascending order (lowest to highest)
	sort.Ints(sorted)

	return sorted
}

// calculateResponseTimePercentiles calculates p50, p95, and p99 for response times
// This gives us a better understanding of response time distribution
func calculateResponseTimePercentiles(responseTimes []int) (p50, p95, p99 int) {
	if len(responseTimes) == 0 {
		return 0, 0, 0
	}

	// Sort the response times
	sorted := sortInts(responseTimes)

	// Calculate percentiles
	p50 = calculatePercentile(sorted, 50) // Median - half of requests are faster than this
	p95 = calculatePercentile(sorted, 95) // 95% of requests are faster than this
	p99 = calculatePercentile(sorted, 99) // 99% of requests are faster than this

	return p50, p95, p99
}

// Why percentiles matter:
//
// Average can be misleading because a few very slow requests can skew it.
// Example:
//   - 98 requests take 100ms
//   - 2 requests take 5000ms (5 seconds)
//   - Average: (98*100 + 2*5000) / 100 = 198ms
//   - But most requests (98%) are actually 100ms!
//
// Percentiles show the full picture:
//   - p50 (median): 100ms - half of users experience this or better
//   - p95: 100ms - 95% of users experience this or better
//   - p99: 5000ms - only 1% of users have this slow experience
//
// This helps you understand:
//   - p50: Typical user experience
//   - p95: Most users' experience (excludes worst 5%)
//   - p99: Worst-case scenarios (excludes only the absolute worst 1%)
