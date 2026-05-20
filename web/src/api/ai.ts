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
      onError(`HTTP error: ${response.status}`)
      return
    }

    const reader = response.body?.getReader()
    if (!reader) {
      onError('No response body')
      return
    }

    const decoder = new TextDecoder()
    while (true) {
      const { done, value } = await reader.read()
      if (done) break
      const text = decoder.decode(value, { stream: true })
      onMessage(text)
    }
    onDone()
  } catch (error: any) {
    onError(error.message || 'Unknown error')
  }
}
