import request from './index'

export const getDashboard = () => request.get('/admin/dashboard')
export const getServerMonitor = () => request.get('/admin/monitor')
export const getVisitMapData = () => request.get('/admin/visit-map')

// 文章管理
export const adminGetArticleList = (params: Record<string, any>) => request.get('/admin/article/list', { params })
export const createArticle = (data: Record<string, any>) => request.post('/admin/article', data)
export const adminGetArticle = (id: number) => request.get(`/admin/article/${id}`)
export const updateArticle = (id: number, data: Record<string, any>) => request.put(`/admin/article/${id}`, data)
export const deleteArticle = (id: number) => request.delete(`/admin/article/${id}`)
export const updateArticleStatus = (id: number, data: Record<string, any>) => request.put(`/admin/article/${id}/status`, data)

// 分类管理
export const adminGetCategories = () => request.get('/admin/category/list')
export const createCategory = (data: Record<string, any>) => request.post('/admin/category', data)
export const updateCategory = (id: number, data: Record<string, any>) => request.put(`/admin/category/${id}`, data)
export const deleteCategory = (id: number) => request.delete(`/admin/category/${id}`)

// 标签管理
export const adminGetTags = () => request.get('/admin/tag/list')
export const createTag = (data: Record<string, any>) => request.post('/admin/tag', data)
export const updateTag = (id: number, data: Record<string, any>) => request.put(`/admin/tag/${id}`, data)
export const deleteTag = (id: number) => request.delete(`/admin/tag/${id}`)

// 评论管理
export const adminGetComments = (params: Record<string, any>) => request.get('/admin/comment/list', { params })
export const reviewComment = (id: number, data: Record<string, any>) => request.put(`/admin/comment/${id}/review`, data)
export const deleteComment = (id: number) => request.delete(`/admin/comment/${id}`)
export const batchReviewComments = (data: Record<string, any>) => request.post('/admin/comment/batch-review', data)

// 用户管理
export const adminGetUsers = (params: Record<string, any>) => request.get('/admin/user/list', { params })
export const adminGetUser = (id: number) => request.get(`/admin/user/${id}`)
export const updateUserStatus = (id: number, data: Record<string, any>) => request.put(`/admin/user/${id}/status`, data)
export const resetUserPassword = (id: number) => request.put(`/admin/user/${id}/reset-password`)

// 管理员管理
export const getAdminList = () => request.get('/admin/admin-mgmt/list')
export const createAdmin = (data: Record<string, any>) => request.post('/admin/admin-mgmt', data)
export const updateAdmin = (id: number, data: Record<string, any>) => request.put(`/admin/admin-mgmt/${id}`, data)
export const updateAdminStatus = (id: number, data: Record<string, any>) => request.put(`/admin/admin-mgmt/${id}/status`, data)

// 系统配置
export const getSystemConfig = () => request.get('/admin/config')
export const updateSystemConfig = (data: Record<string, any>) => request.put('/admin/config', data)
export const getSiteConfig = () => request.get('/site/config')

// AI模型管理
export const getAIModelList = () => request.get('/admin/ai-model/list')
export const createAIModel = (data: Record<string, any>) => request.post('/admin/ai-model', data)
export const updateAIModel = (id: number, data: Record<string, any>) => request.put(`/admin/ai-model/${id}`, data)
export const updateAIModelStatus = (id: number, data: Record<string, any>) => request.put(`/admin/ai-model/${id}/status`, data)
export const deleteAIModel = (id: number) => request.delete(`/admin/ai-model/${id}`)
export const updateAIModelGlobalConfig = (data: Record<string, any>) => request.put('/admin/ai-model/global-config', data)

// 备份
export const getBackupList = () => request.get('/admin/backup/list')
export const createBackup = () => request.post('/admin/backup')
export const downloadBackup = (id: number) => request.get(`/admin/backup/${id}/download`, { responseType: 'blob' })
export const deleteBackup = (id: number) => request.delete(`/admin/backup/${id}`)
export const updateAutoBackupConfig = (data: Record<string, any>) => request.put('/admin/backup/auto-config', data)

// 操作日志
export const getOperationLogs = (params: Record<string, any>) => request.get('/admin/logs', { params })

// 图片上传
export const uploadImage = (formData: FormData) => request.post('/admin/upload', formData)
