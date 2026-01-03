// formatters.js - Utility functions for formatting data

export const formatTimestamp = (timestamp) => {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('en-US', {
    hour12: false,
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    fractionalSecondDigits: 3
  })
}

export const formatDuration = (seconds) => {
  if (seconds < 60) {
    return `${seconds.toFixed(2)}s`
  }
  const minutes = Math.floor(seconds / 60)
  const remainingSeconds = seconds % 60
  return `${minutes}m ${remainingSeconds.toFixed(0)}s`
}

export const formatStatusColor = (status) => {
  if (status >= 200 && status < 300) return 'text-green-600 bg-green-50'
  if (status >= 400 && status < 500) return 'text-yellow-600 bg-yellow-50'
  if (status >= 500) return 'text-red-600 bg-red-50'
  return 'text-gray-600 bg-gray-50'
}

export const formatStatusText = (status) => {
  const statusCodes = {
    200: 'OK',
    201: 'Created',
    204: 'No Content',
    400: 'Bad Request',
    401: 'Unauthorized',
    403: 'Forbidden',
    404: 'Not Found',
    406: 'Not Acceptable',
    429: 'Too Many Requests',
    500: 'Internal Server Error',
    502: 'Bad Gateway',
    503: 'Service Unavailable',
    504: 'Gateway Timeout'
  }
  return statusCodes[status] || `Status ${status}`
}

export const truncateText = (text, maxLength = 100) => {
  if (!text) return ''
  if (text.length <= maxLength) return text
  return text.substring(0, maxLength) + '...'
}

export const calculatePercentage = (value, total) => {
  if (total === 0) return 0
  return ((value / total) * 100).toFixed(1)
}
