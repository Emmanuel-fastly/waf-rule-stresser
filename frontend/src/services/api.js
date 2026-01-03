// api.js - API service for backend communication

const API_BASE_URL = '/api'

export const startTest = async (config) => {
  const response = await fetch(`${API_BASE_URL}/test/start`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(config),
  })

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Test failed: ${errorText}`)
  }

  return response.json()
}

export const exportResults = async (config, results, format) => {
  const response = await fetch(`${API_BASE_URL}/test/export`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ config, results, format }),
  })

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Export failed: ${errorText}`)
  }

  return response.json()
}

export const listExports = async () => {
  const response = await fetch(`${API_BASE_URL}/exports/list`)

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Failed to list exports: ${errorText}`)
  }

  return response.json()
}

// Stream a test with real-time progress updates via SSE
export const startTestStreaming = (config, onProgress, onComplete, onError) => {
  // Make initial POST request to start the test
  fetch(`${API_BASE_URL}/test/stream`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(config),
  })
    .then(response => {
      if (!response.ok) {
        throw new Error(`Failed to start test: ${response.statusText}`)
      }

      // Set up EventSource to receive SSE messages
      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      // Read stream
      const readStream = () => {
        reader.read().then(({ done, value }) => {
          if (done) {
            return
          }

          // Decode chunk
          buffer += decoder.decode(value, { stream: true })

          // Process complete SSE messages
          const lines = buffer.split('\n\n')
          buffer = lines.pop() // Keep incomplete message in buffer

          lines.forEach(line => {
            if (line.startsWith('data: ')) {
              try {
                const jsonData = line.substring(6) // Remove 'data: ' prefix
                const update = JSON.parse(jsonData)

                switch (update.type) {
                  case 'progress':
                    onProgress(update)
                    break

                  case 'complete':
                    onComplete(update.final_result)
                    break

                  case 'error':
                    onError(update.error)
                    break

                  case 'cancelled':
                    onComplete(null) // Signal cancellation
                    break
                }
              } catch (err) {
                console.error('Failed to parse SSE message:', err)
              }
            }
          })

          // Continue reading
          readStream()
        }).catch(err => {
          console.error('Stream reading error:', err)
          onError('Connection to server lost')
        })
      }

      readStream()
    })
    .catch(err => {
      console.error('Failed to start streaming test:', err)
      onError(err.message)
    })
}

// Stop a running test
export const stopTest = async (testID) => {
  const response = await fetch(`${API_BASE_URL}/test/stop`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ test_id: testID }),
  })

  if (!response.ok) {
    const errorText = await response.text()
    throw new Error(`Failed to stop test: ${errorText}`)
  }

  return response.json()
}
