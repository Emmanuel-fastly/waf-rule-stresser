import { formatDuration } from '../utils/formatters'

const ResponseTimeStats = ({ results }) => {
  if (!results) return null

  const stats = [
    { label: 'Avg Response', value: results.avg_response, unit: 'ms' },
    { label: 'Min Response', value: results.min_response, unit: 'ms' },
    { label: 'Max Response', value: results.max_response, unit: 'ms' },
    { label: 'P50 (Median)', value: results.p50_response, unit: 'ms' },
    { label: 'P95', value: results.p95_response, unit: 'ms' },
    { label: 'P99', value: results.p99_response, unit: 'ms' },
    { label: 'Duration', value: formatDuration(results.duration), unit: '' },
    { label: 'Requests/Sec', value: results.requests_per_sec?.toFixed(2), unit: '' }
  ]

  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <h2 className="text-xl font-bold text-gray-800 mb-4">Performance Metrics</h2>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {stats.map((stat, index) => (
          <div key={index} className="border border-gray-200 rounded-lg p-4">
            <div className="text-2xl font-bold text-gray-900">
              {stat.value}{stat.unit && <span className="text-sm text-gray-500 ml-1">{stat.unit}</span>}
            </div>
            <div className="text-xs text-gray-600 mt-1">{stat.label}</div>
          </div>
        ))}
      </div>

      {results.rate_limit_hit && (
        <div className="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
          <div className="flex items-center gap-2">
            <span className="text-xl">🚨</span>
            <div>
              <div className="font-semibold text-red-800">Rate Limit Detected</div>
              <div className="text-sm text-red-600">
                First blocked at request #{results.rate_limit_at}
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default ResponseTimeStats
