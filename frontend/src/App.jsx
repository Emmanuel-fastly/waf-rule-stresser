import { useState, useRef, useEffect } from 'react'
import Header from './components/Header'
import TestConfiguration from './components/TestConfiguration'
import StatsCards from './components/StatsCards'
import ResponseTimeStats from './components/ResponseTimeStats'
import RequestLog from './components/RequestLog'
import ProgressBar from './components/ProgressBar'
import { startTestStreaming, stopTest, exportResults } from './services/api'

const App = () => {
  const [testResults, setTestResults] = useState(null)
  const [testConfig, setTestConfig] = useState(null)
  const [error, setError] = useState(null)
  const [exportStatus, setExportStatus] = useState(null)

  // Streaming test state
  const [isTestRunning, setIsTestRunning] = useState(false)
  const [progress, setProgress] = useState({ completed: 0, total: 0, percentage: 0 })
  const [liveRequests, setLiveRequests] = useState([])
  const [liveStats, setLiveStats] = useState(null)
  const [currentTestId, setCurrentTestId] = useState(null)

  // Refs for one-time scroll and tracking live data
  const liveRequestsRef = useRef([])
  const liveStatsRef = useRef(null)
  const hasScrolledToLogsRef = useRef(false)
  const requestLogRef = useRef(null)

  // One-time scroll when request log appears
  useEffect(() => {
    if (isTestRunning && liveRequests.length > 0 && !hasScrolledToLogsRef.current && requestLogRef.current) {
      // Scroll to bring request log into view smoothly
      requestLogRef.current.scrollIntoView({ behavior: 'smooth', block: 'start' })
      hasScrolledToLogsRef.current = true
    }
  }, [isTestRunning, liveRequests.length])

  // Streaming test handler
  const handleStartTestStreaming = (config) => {
    // Reset state
    setIsTestRunning(true)
    setError(null)
    setTestResults(null)
    setTestConfig(config)
    setExportStatus(null)
    setProgress({ completed: 0, total: config.total_requests, percentage: 0 })
    setLiveRequests([])
    setLiveStats(null)

    // Reset refs
    liveRequestsRef.current = []
    liveStatsRef.current = null
    hasScrolledToLogsRef.current = false

    // Start streaming
    startTestStreaming(
      config,
      // onProgress callback
      (update) => {
        setCurrentTestId(update.test_id)
        setProgress({
          completed: update.completed,
          total: update.total,
          percentage: update.percentage
        })

        // Append new requests to live log and update ref
        if (update.new_requests && update.new_requests.length > 0) {
          setLiveRequests(prev => {
            const updated = [...prev, ...update.new_requests]
            liveRequestsRef.current = updated
            return updated
          })
        }

        // Update live stats and ref
        if (update.current_stats) {
          setLiveStats(update.current_stats)
          liveStatsRef.current = update.current_stats
        }
      },
      // onComplete callback
      (finalResult) => {
        setIsTestRunning(false)
        if (finalResult) {
          setTestResults(finalResult)
        } else {
          // Test was cancelled - use refs to get latest values
          // This fixes the closure issue where state might be stale
          const currentRequests = liveRequestsRef.current
          const currentStats = liveStatsRef.current

          if (currentRequests.length > 0) {
            const partialResult = {
              test_id: currentTestId,
              total_requests: currentRequests.length,
              success_count: currentStats?.success_count || 0,
              error_count: currentStats?.error_count || 0,
              blocked_count: currentStats?.blocked_count || 0,
              avg_response: currentStats?.avg_response || 0,
              min_response: 0,
              max_response: 0,
              p50_response: 0,
              p95_response: 0,
              p99_response: 0,
              requests: currentRequests,
              rate_limit_hit: false,
              rate_limit_at: 0,
              start_time: new Date().toISOString(),
              end_time: new Date().toISOString(),
              duration: 0,
              requests_per_sec: 0
            }
            setTestResults(partialResult)
          }
        }
      },
      // onError callback
      (errorMessage) => {
        setIsTestRunning(false)
        setError(errorMessage)
      }
    )
  }

  // Stop test handler
  const handleStopTest = async () => {
    if (!currentTestId) return

    try {
      // Send stop request to backend
      await stopTest(currentTestId)

      setIsTestRunning(false)
      // Don't show error message for user-initiated stops
    } catch (err) {
      console.error('Failed to stop test:', err)
      setError(`Failed to stop test: ${err.message}`)
    }
  }

  const handleExport = async (format) => {
    if (!testResults || !testConfig) return

    try {
      setExportStatus('Exporting...')
      const response = await exportResults(testConfig, testResults, format)
      setExportStatus(`✓ Exported to ${response.filename}`)
      setTimeout(() => setExportStatus(null), 5000)
    } catch (err) {
      setExportStatus(`✗ Export failed: ${err.message}`)
      setTimeout(() => setExportStatus(null), 5000)
    }
  }

  const handleClearLogs = () => {
    setTestResults(null)
    setTestConfig(null)
    setError(null)
    setExportStatus(null)
    setLiveRequests([])
    setLiveStats(null)
    setProgress({ completed: 0, total: 0, percentage: 0 })
  }

  return (
    <div className="min-h-screen bg-gray-100">
      <Header />

      <main className="max-w-7xl mx-auto px-4 py-8 space-y-8">
        {/* Test Configuration */}
        <TestConfiguration
          onStartTest={handleStartTestStreaming}
          onStopTest={handleStopTest}
          onClearLogs={handleClearLogs}
          isTestRunning={isTestRunning}
        />

        {/* Progress Bar - shown during test */}
        {isTestRunning && (
          <ProgressBar
            progress={progress}
            liveStats={liveStats}
          />
        )}

        {/* Error State */}
        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-6">
            <div className="flex items-start gap-3">
              <span className="text-2xl">❌</span>
              <div>
                <h3 className="font-semibold text-red-800 mb-1">Test Failed</h3>
                <p className="text-red-700 text-sm">{error}</p>
              </div>
            </div>
          </div>
        )}

        {/* Live Request Log (during test) */}
        {isTestRunning && liveRequests.length > 0 && (
          <div ref={requestLogRef}>
            <RequestLog
              results={{ requests: liveRequests }}
              isLive={true}
            />
          </div>
        )}

        {/* Results */}
        {testResults && !isTestRunning && (
          <>
            {/* Export Controls */}
            <div className="bg-white rounded-lg shadow-md p-6">
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
                <div>
                  <h2 className="text-xl font-bold text-gray-800">Test Results</h2>
                  <p className="text-sm text-gray-600 mt-1">
                    Test ID: {testResults.test_id}
                  </p>
                </div>
                <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-2 sm:gap-3">
                  {exportStatus && (
                    <span className="text-sm text-gray-700">{exportStatus}</span>
                  )}
                  <button
                    onClick={() => handleExport('json')}
                    className="px-3 sm:px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition text-sm font-medium"
                  >
                    Export JSON
                  </button>
                  <button
                    onClick={() => handleExport('csv')}
                    className="px-3 sm:px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition text-sm font-medium"
                  >
                    Export CSV
                  </button>
                </div>
              </div>
            </div>

            {/* Stats Cards */}
            <StatsCards results={testResults} />

            {/* Performance Metrics */}
            <ResponseTimeStats results={testResults} />

            {/* Request Log */}
            <RequestLog results={testResults} />
          </>
        )}

        {/* Empty State */}
        {!isTestRunning && !testResults && !error && (
          <div className="bg-white rounded-lg shadow-md p-12 text-center">
            <h2 className="text-2xl font-bold text-gray-800 mb-2">
              Ready to Test
            </h2>
            <p className="text-gray-600">
              Configure your test parameters above and click "Start Test" to begin
            </p>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="bg-white border-t border-gray-200 mt-12">
        <div className="max-w-7xl mx-auto px-4 py-6 text-center text-sm text-gray-600">
          <p>WAF Block & Rate Limit Tester - Built with Go + React + Tailwind</p>
        </div>
      </footer>
    </div>
  )
}

export default App
