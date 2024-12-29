# 阶段 1: 构建后端
FROM golang:1.20-alpine as backend-builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-visitor

# 阶段 2: 运行环境
FROM nginx:alpine
WORKDIR /app

# 安装 supervisord
RUN apk add --no-cache supervisor

# 直接复制已构建的前端文件
COPY front/dist /usr/share/nginx/html

# 直接复制已构建的后端二进制文件
COPY ./flash_visit /app/flash_visit

# 配置 nginx
COPY nginx.conf /etc/nginx/conf.d/default.conf

# 配置 supervisord
COPY supervisord.conf /etc/supervisord.conf

# 创建日志目录
RUN mkdir -p /var/log/supervisor

# 暴露端口
EXPOSE 3000

# 启动服务
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
