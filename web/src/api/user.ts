import request from './index'

export const getUserInfo = () => request.get('/user/info')
export const updateUserInfo = (data: Record<string, any>) => request.put('/user/info', data)
export const changePassword = (data: Record<string, any>) => request.put('/user/password', data)
export const uploadAvatar = (formData: FormData) => request.post('/user/avatar', formData)
export const getSubscription = () => request.get('/user/subscription')
export const updateSubscription = (data: Record<string, any>) => request.put('/user/subscription', data)
export const getAIConfig = () => request.get('/user/ai-config')
export const updateAIConfig = (data: Record<string, any>) => request.put('/user/ai-config', data)
export const getAvailableModels = () => request.get('/user/ai-models')
