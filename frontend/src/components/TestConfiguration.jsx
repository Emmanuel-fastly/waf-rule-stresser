import { useState } from 'react'
import { Play, Pause, Shield, Zap, AlertTriangle, Clock, Gauge } from 'lucide-react'

const TestConfiguration = ({ onStartTest, onStopTest, onClearLogs, isTestRunning }) => {
  const [config, setConfig] = useState({
    target_url: '',
    total_requests: 100,
    duration: 10,
    traffic_type: 'normal',
    error_mode: false,
    user_agent_type: 'legitimate',
    custom_user_agent: '',
    test_mode: 'baseline',
    http_method: 'GET',
    custom_headers: {},
    request_body: ''
  })

  const [headerKey, setHeaderKey] = useState('')
  const [headerValue, setHeaderValue] = useState('')

  const handleSubmit = (e) => {
    e.preventDefault()

    // If test is running, stop it
    if (isTestRunning) {
      onStopTest()
      return
    }

    // Validate URL
    if (!config.target_url.trim()) {
      alert('Please enter a target URL')
      return
    }

    // Check if URL has a valid scheme (http:// or https://)
    if (!config.target_url.match(/^https?:\/\//i)) {
      alert('URL must start with http:// or https://')
      return
    }

    onStartTest(config)
  }

  const handleChange = (field, value) => {
    setConfig(prev => ({ ...prev, [field]: value }))
  }

  const addHeader = () => {
    if (headerKey && headerValue) {
      setConfig(prev => ({
        ...prev,
        custom_headers: { ...prev.custom_headers, [headerKey]: headerValue }
      }))
      setHeaderKey('')
      setHeaderValue('')
    }
  }

  const removeHeader = (key) => {
    setConfig(prev => {
      const newHeaders = { ...prev.custom_headers }
      delete newHeaders[key]
      return { ...prev, custom_headers: newHeaders }
    })
  }

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold text-gray-800 mb-6">Test Configuration</h2>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Target URL and Traffic Type */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Target URL */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Target URL
            </label>
            <input
              type="text"
              value={config.target_url}
              onChange={(e) => handleChange('target_url', e.target.value)}
              placeholder="https://httpbin.org/get"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          {/* Traffic Type */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Traffic Type
            </label>
            <div className="grid grid-cols-2 gap-3">
              <button
                type="button"
                onClick={() => handleChange('traffic_type', 'normal')}
                className={`px-6 py-2 rounded-md font-medium transition flex items-center justify-center gap-2 ${
                  config.traffic_type === 'normal'
                    ? 'bg-blue-600 text-white'
                    : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 hover:border-blue-500'
                }`}
              >
                <Shield size={18} />
                Normal
              </button>
              <button
                type="button"
                onClick={() => handleChange('traffic_type', 'attack')}
                className={`px-6 py-2 rounded-md font-medium transition flex items-center justify-center gap-2 ${
                  config.traffic_type === 'attack'
                    ? 'bg-red-600 text-white'
                    : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 hover:border-red-500'
                }`}
              >
                <Zap size={18} />
                Attack
              </button>
            </div>
            {/* Error Mode checkbox - disabled for attack traffic */}
            <div className="mt-3 flex items-center">
              <input
                type="checkbox"
                id="error_mode"
                checked={config.error_mode}
                onChange={(e) => handleChange('error_mode', e.target.checked)}
                disabled={config.traffic_type === 'attack'}
                className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded disabled:opacity-50 disabled:cursor-not-allowed"
              />
              <label
                htmlFor="error_mode"
                className={`ml-2 block text-sm ${
                  config.traffic_type === 'attack' ? 'text-gray-400' : 'text-gray-700'
                }`}
              >
                404 (Force 404 responses)
              </label>
            </div>
          </div>
        </div>

        {/* Request Configuration */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Total Requests
            </label>
            <input
              type="number"
              value={config.total_requests}
              onChange={(e) => handleChange('total_requests', parseInt(e.target.value))}
              min="1"
              max="10000"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Duration (seconds)
            </label>
            <input
              type="number"
              value={config.duration}
              onChange={(e) => handleChange('duration', parseInt(e.target.value))}
              min="1"
              max="3600"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              HTTP Method
            </label>
            <select
              value={config.http_method}
              onChange={(e) => handleChange('http_method', e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
              <option value="PATCH">PATCH</option>
            </select>
          </div>
        </div>


        {/* User Agent Configuration */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              User Agent Type
            </label>
            <select
              value={config.user_agent_type}
              onChange={(e) => handleChange('user_agent_type', e.target.value)}
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none"
              disabled={config.custom_user_agent !== ''}
            >
              <option value="legitimate">Legitimate Browser</option>
              <option value="scanner">Scanner/Bot</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Custom User Agent (Optional)
            </label>
            <input
              type="text"
              value={config.custom_user_agent}
              onChange={(e) => handleChange('custom_user_agent', e.target.value)}
              placeholder="Custom-Agent/1.0"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>


        {/* Request Body (for POST/PUT/PATCH) */}
        {['POST', 'PUT', 'PATCH'].includes(config.http_method) && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Request Body (JSON)
            </label>
            <textarea
              value={config.request_body}
              onChange={(e) => handleChange('request_body', e.target.value)}
              placeholder='{"key": "value"}'
              rows="3"
              className="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
            />
          </div>
        )}

        {/* Custom Headers */}
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-2">
            Custom Headers
          </label>
          <div className="flex flex-col sm:flex-row gap-2 mb-2">
            <input
              type="text"
              value={headerKey}
              onChange={(e) => setHeaderKey(e.target.value)}
              placeholder="Header-Name"
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
            <input
              type="text"
              value={headerValue}
              onChange={(e) => setHeaderValue(e.target.value)}
              placeholder="Header-Value"
              className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
            <button
              type="button"
              onClick={addHeader}
              className="px-8 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition whitespace-nowrap"
            >
              Add
            </button>
          </div>
          {Object.keys(config.custom_headers).length > 0 && (
            <div className="space-y-1">
              {Object.entries(config.custom_headers).map(([key, value]) => (
                <div key={key} className="flex items-center justify-between bg-gray-50 px-3 py-2 rounded">
                  <span className="text-sm font-mono">
                    <span className="font-semibold">{key}:</span> {value}
                  </span>
                  <button
                    type="button"
                    onClick={() => removeHeader(key)}
                    className="text-red-600 hover:text-red-800 text-sm"
                  >
                    Remove
                  </button>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Bottom Section: Action Buttons and Test Mode */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Start Test and Clear Logs Buttons */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2 invisible">Actions</label>
            <div className="grid grid-cols-2 gap-3">
              <button
                type="submit"
                className={`py-2 px-4 sm:px-6 rounded-md transition font-medium flex items-center justify-center gap-2 text-sm sm:text-base ${
                  isTestRunning
                    ? 'bg-red-600 text-white hover:bg-red-700'
                    : 'bg-blue-600 text-white hover:bg-blue-700'
                }`}
              >
                {isTestRunning ? (
                  <>
                    <Pause size={18} />
                    <span className="hidden sm:inline">Stop Test</span>
                    <span className="sm:hidden">Stop</span>
                  </>
                ) : (
                  <>
                    <Play size={18} />
                    <span className="hidden sm:inline">Start Test</span>
                    <span className="sm:hidden">Start</span>
                  </>
                )}
              </button>
              <button
                type="button"
                onClick={onClearLogs}
                disabled={isTestRunning}
                className="bg-white text-gray-700 py-2 px-4 sm:px-6 rounded-md border border-gray-300 hover:bg-gray-50 transition font-medium text-sm sm:text-base disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Clear Logs
              </button>
            </div>
          </div>

          {/* Test Mode */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Test Mode
            </label>
            <div className="grid grid-cols-2 gap-3">
              <div>
                <button
                  type="button"
                  onClick={() => handleChange('test_mode', 'baseline')}
                  className={`w-full px-4 sm:px-6 py-2 rounded-md font-medium transition flex items-center justify-center gap-2 text-sm sm:text-base ${
                    config.test_mode === 'baseline'
                      ? 'bg-green-600 text-white'
                      : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 hover:border-green-500'
                  }`}
                >
                  <Clock size={18} />
                  Baseline
                </button>
                <p className="text-xs text-gray-500 mt-2">Evenly distributed</p>
              </div>
              <div>
                <button
                  type="button"
                  onClick={() => handleChange('test_mode', 'burst')}
                  className={`w-full px-4 sm:px-6 py-2 rounded-md font-medium transition flex items-center justify-center gap-2 text-sm sm:text-base ${
                    config.test_mode === 'burst'
                      ? 'bg-orange-600 text-white'
                      : 'bg-white text-gray-700 border border-gray-300 hover:bg-gray-50 hover:border-orange-500'
                  }`}
                >
                  <Gauge size={18} />
                  Burst
                </button>
                <p className="text-xs text-gray-500 mt-2">Spike at start</p>
              </div>
            </div>
          </div>
        </div>
      </form>
    </div>
  )
}

export default TestConfiguration
