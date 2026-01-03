import { Activity, CheckCircle, XCircle, ShieldX } from 'lucide-react'

const ProgressBar = ({ progress, liveStats }) => {
  return (
    <div className="bg-white rounded-lg shadow-md p-6">
      <div className="mb-4">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-lg font-semibold text-gray-800">Test Progress</h3>
          <span className="text-2xl font-bold text-blue-600">
            {progress.percentage}%
          </span>
        </div>

        {/* Progress Bar */}
        <div className="w-full bg-gray-200 rounded-full h-4 overflow-hidden">
          <div
            className="bg-blue-600 h-full transition-all duration-300 ease-out flex items-center justify-end pr-2"
            style={{ width: `${progress.percentage}%` }}
          >
            {progress.percentage > 10 && (
              <span className="text-xs font-semibold text-white">
                {progress.completed}/{progress.total}
              </span>
            )}
          </div>
        </div>

        <p className="text-sm text-gray-600 mt-2">
          {progress.completed} of {progress.total} requests completed
        </p>
      </div>

      {/* Live Stats */}
      {liveStats && (
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 pt-4 border-t border-gray-200">
          <div className="flex items-center gap-2">
            <CheckCircle size={16} className="text-green-600" />
            <div>
              <p className="text-xs text-gray-600">Success</p>
              <p className="text-lg font-semibold text-gray-900">{liveStats.success_count}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <XCircle size={16} className="text-red-600" />
            <div>
              <p className="text-xs text-gray-600">Errors</p>
              <p className="text-lg font-semibold text-gray-900">{liveStats.error_count}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <ShieldX size={16} className="text-purple-600" />
            <div>
              <p className="text-xs text-gray-600">Blocked</p>
              <p className="text-lg font-semibold text-gray-900">{liveStats.blocked_count}</p>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Activity size={16} className="text-blue-600" />
            <div>
              <p className="text-xs text-gray-600">Avg Response</p>
              <p className="text-lg font-semibold text-gray-900">{liveStats.avg_response}ms</p>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default ProgressBar
