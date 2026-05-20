# GoBlog - 个人博客系统

## 项目简介
GoBlog 是一款基于 Go + Vue3 开发的全功能个人博客系统，支持 Markdown 编辑、文章管理、用户交互、后台监控、AI 对话、消息推送等核心功能。

## 技术栈
### 后端
- Go 1.22 + Gin
- GORM + MySQL 8.0
- JWT 认证
- Viper 配置管理

### 前端
- Vue 3 + TypeScript
- Ant Design Vue + Ant Design X Vue
- Vue I18n（中英双语）
- Pinia 状态管理
- ECharts 图表
- Markdown-it + Highlight.js

## 功能特性
### 前台功能
- 用户注册/登录（邮箱验证码、JWT鉴权）
- 文章浏览（列表、详情、搜索、分类/标签筛选）
- 文章交互（浏览量、点赞、评论）
- AI文章摘要（基于用户配置的AI模型）
- AI对话助手（流式输出、文章上下文）
- 个人中心（信息修改、订阅配置、AI配置）
- RSS订阅
- 深浅色主题切换
- 中英文切换
- 响应式设计

### 后台功能
- 数据仪表盘（访问量、文章数、字数统计）
- 服务器实时监控（CPU/内存/磁盘/网络）
- 访客地域可视化（2D/3D地图、飞线动画）
- 文章管理（CRUD、Markdown编辑器、图片上传）
- 评论管理（审核、批量操作）
- 用户管理（查看、禁用、重置密码）
- 管理员管理（仅超级管理员）
- 系统配置（站点信息、注册开关、SMTP、飞书）
- AI模型管理（模型池、全局控制）
- 数据备份（手动/自动）
- 操作日志

## 快速开始

### 环境要求
- Go 1.22+
- Node.js 20+
- MySQL 8.0+
- Docker & Docker Compose（可选）

### 方式一：Docker 部署（推荐）
```bash
# 克隆项目
git clone https://github.com/yourname/GoBlog.git
cd GoBlog

# 启动服务
docker-compose up -d

# 访问
# 前台：http://localhost:8080
# 后台：http://localhost:8080/admin
# 默认管理员：admin / admin123
```

### 方式二：本地开发
```bash
# 后端
cd GoBlog
go mod tidy
# 修改 config.yaml 中的数据库配置
go run main.go

# 前端
cd web
npm install
npm run dev
```

### 配置说明
编辑 `config.yaml` 文件进行配置：
- 数据库连接信息
- JWT密钥（生产环境务必修改）
- SMTP邮件配置
- 飞书机器人配置
- 上传文件配置
- AI功能开关
- 监控开关

## 项目结构
```
GoBlog/
├── config/          # 配置
├── database/        # 数据库初始化
├── handler/         # HTTP处理器
├── middleware/       # 中间件
├── model/           # 数据模型
├── router/          # 路由
├── service/         # 业务服务
├── utils/           # 工具函数
├── web/             # 前端项目
│   ├── src/
│   │   ├── api/     # API请求
│   │   ├── assets/  # 静态资源
│   │   ├── components/ # 公共组件
│   │   ├── composables/ # 组合式函数
│   │   ├── i18n/    # 国际化
│   │   ├── layouts/ # 布局
│   │   ├── router/  # 路由
│   │   ├── stores/  # 状态管理
│   │   └── views/   # 页面
│   └── ...
├── config.yaml      # 配置文件
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## API 接口
### 公开接口
- GET /api/site/config - 获取站点配置
- GET /api/article/list - 文章列表
- GET /api/article/:id - 文章详情
- POST /api/article/:id/like - 点赞
- GET /api/article/:id/comments - 评论列表
- POST /api/article/:id/comment - 提交评论
- GET /api/categories - 分类列表
- GET /api/tags - 标签列表
- GET /api/rss - RSS订阅

### 用户接口（需登录）
- POST /api/auth/register - 注册
- POST /api/auth/login - 登录
- GET /api/user/info - 用户信息
- PUT /api/user/info - 更新信息
- POST /api/user/chat - AI对话

### 管理员接口（需管理员登录）
- GET /api/admin/dashboard - 仪表盘数据
- GET /api/admin/monitor - 服务器监控
- POST /api/admin/article - 创建文章
- GET /api/admin/comment/list - 评论列表
- ...

## 许可证
MIT License
