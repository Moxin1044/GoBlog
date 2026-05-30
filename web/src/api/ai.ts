import request from './index'

export const chat = (data: { message: string; article_id?: number }) => request.post('/user/chat', data)
export const getChatHistory = (params?: Record<string, any>) => request.get('/user/chat/history', { params })

export const streamChat = async (
  data: { message: string; article_id?: number },
  onMessage: (content: string) => void,
  onDone: () => void,
  onError: (error: string) => void
) => {
  const token = localStorage.getItem('token')
  try {
    const response = await fetch('/api/user/chat', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
      },
      body: JSON.stringify(data),
    })

    if (!response.ok) {
      const errText = await response.text()
      try {
        const errJson = JSON.parse(errText)
        onError(errJson.message || `HTTP error: ${response.status}`)
      } catch {
        onError(`HTTP error: ${response.status}`)
      }
      return
    }

    const reader = response.body?.getReader()
    if (!reader) {
      onError('No response body')
      return
    }

    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        const trimmed = line.trim()
        if (!trimmed || trimmed.startsWith(':')) continue

        if (trimmed.startsWith('event:')) continue

        if (trimmed.startsWith('data:')) {
          const dataStr = trimmed.slice(5).trim()
          if (dataStr === '[DONE]') {
            onDone()
            return
          }

          try {
            const parsed = JSON.parse(dataStr)
            if (parsed.done) {
              onDone()
              return
            }
            if (parsed.content) {
              onMessage(parsed.content)
            }
          } catch {
            // skip malformed JSON
          }
        }
      }
    }

    onDone()
  } catch (error: any) {
    onError(error.message || 'Unknown error')
  }
}
