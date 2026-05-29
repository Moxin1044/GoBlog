import request from './index'

export const login = (data: { account: string; password: string }) => request.post('/auth/login', data)
export const register = (data: Record<string, any>) => request.post('/auth/register', data)
export const sendVerifyCode = (email: string, type: string = 'register') => request.post('/auth/verify-code', { email, type })
