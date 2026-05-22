import request from './index'

export const getArticleList = (params: Record<string, any>) => request.get('/article/list', { params })
export const getArticleDetail = (id: number) => request.get(`/article/${id}`)
export const likeArticle = (id: number) => request.post(`/article/${id}/like`)
export const getComments = (id: number, params?: Record<string, any>) => request.get(`/article/${id}/comments`, { params })
export const submitComment = (id: number, data: Record<string, any>) => request.post(`/article/${id}/comment`, data)
export const getCategories = () => request.get('/categories')
export const getTags = () => request.get('/tags')
export const getNavigations = () => request.get('/navigations')
export const generateAISummary = (id: number) => request.post(`/article/${id}/ai-summary`)
