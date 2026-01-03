import { formatTimestamp, truncateText } from '../utils/formatters'

const RequestDetails = ({ request, onClose }) => {
  if (!request) return null

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
          <h2 className="text-xl font-bold text-gray-800">
            Request #{request.id} Details
          </h2>
          <button
            onClick={onClose}
            className="text-gray-500 hover:text-gray-700 text-2xl font-bold"
          >
            ×
          </button>
        </div>

        <div className="p-6 space-y-6">
          {/* Status */}
          <div>
            <h3 className="text-sm font-semibold text-gray-700 mb-2">Status</h3>
            <div className="flex items-center gap-4">
              <span className={`px-3 py-1 rounded-full text-sm font-semibold ${
                request.status >= 200 && request.status < 300 ? 'bg-green-100 text-green-800' :
                request.status >= 400 && request.status < 500 ? 'bg-yellow-100 text-yellow-800' :
                'bg-red-100 text-red-800'
              }`}>
                {request.status} {request.status_text}
              </span>
              {request.was_blocked && (
                <span className="px-3 py-1 rounded-full text-sm font-semibold bg-red-100 text-red-800">
                  🚫 Blocked
                </span>
              )}
            </div>
          </div>

          {/* Request Info */}
          <div>
            <h3 className="text-sm font-semibold text-gray-700 mb-2">Request Info</h3>
            <div className="bg-gray-50 rounded-lg p-4 space-y-2 font-mono text-sm">
              <div><span className="text-gray-600">Method:</span> <span className="font-semibold">{request.method}</span></div>
              <div><span className="text-gray-600">URL:</span> <span className="break-all">{request.url}</span></div>
              <div><span className="text-gray-600">Time:</span> {formatTimestamp(request.timestamp)}</div>
              <div><span className="text-gray-600">Response Time:</span> <span className="font-semibold">{request.response_time}ms</span></div>
              {request.attack_info && (
                <div><span className="text-gray-600">Attack Type:</span> <span className="text-red-600 font-semibold">{request.attack_info}</span></div>
              )}
            </div>
          </div>

          {/* Request Headers */}
          {request.request_headers && Object.keys(request.request_headers).length > 0 && (
            <div>
              <h3 className="text-sm font-semibold text-gray-700 mb-2">Request Headers</h3>
              <div className="bg-gray-50 rounded-lg p-4 font-mono text-xs space-y-1">
                {Object.entries(request.request_headers).map(([key, value]) => (
                  <div key={key} className="break-all">
                    <span className="text-blue-600 font-semibold">{key}:</span> {value}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Request Body */}
          {request.request_body && (
            <div>
              <h3 className="text-sm font-semibold text-gray-700 mb-2">Request Body</h3>
              <pre className="bg-gray-50 rounded-lg p-4 font-mono text-xs overflow-x-auto">
                {request.request_body}
              </pre>
            </div>
          )}

          {/* Response Headers */}
          {request.response_headers && Object.keys(request.response_headers).length > 0 && (
            <div>
              <h3 className="text-sm font-semibold text-gray-700 mb-2">Response Headers</h3>
              <div className="bg-gray-50 rounded-lg p-4 font-mono text-xs space-y-1 max-h-60 overflow-y-auto">
                {Object.entries(request.response_headers).map(([key, value]) => (
                  <div key={key} className="break-all">
                    <span className="text-green-600 font-semibold">{key}:</span> {value}
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Response Body */}
          {request.response_body && (
            <div>
              <h3 className="text-sm font-semibold text-gray-700 mb-2">Response Body</h3>
              <pre className="bg-gray-50 rounded-lg p-4 font-mono text-xs overflow-x-auto max-h-60 overflow-y-auto">
                {request.response_body}
              </pre>
            </div>
          )}

          {/* Error */}
          {request.error && (
            <div>
              <h3 className="text-sm font-semibold text-gray-700 mb-2">Error</h3>
              <div className="bg-red-50 border border-red-200 rounded-lg p-4 text-sm text-red-800">
                {request.error}
              </div>
            </div>
          )}
        </div>

        <div className="sticky bottom-0 bg-gray-50 px-6 py-4 border-t border-gray-200">
          <button
            onClick={onClose}
            className="w-full bg-gray-600 text-white py-2 px-4 rounded-md hover:bg-gray-700 transition"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  )
}

export default RequestDetails
