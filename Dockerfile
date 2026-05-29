# 阶段1: 构建前端
FROM node:20-alpine AS frontend-builder
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN node node_modules/vite/bin/vite.js build

# 阶段2: 构建后端
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 复制前端构建产物
COPY --from=frontend-builder /app/web/dist ./web/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o goblog .
# 创建运行时需要的目录
RUN mkdir -p uploads backups

# 阶段3: 运行
FROM alpine:3.19
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/goblog .
COPY --from=backend-builder /app/config.yaml .
COPY --from=backend-builder /app/web/dist ./web/dist
COPY --from=backend-builder /app/uploads ./uploads
COPY --from=backend-builder /app/backups ./backups

ENV TZ=Asia/Shanghai
EXPOSE 8080

CMD ["./goblog"]
