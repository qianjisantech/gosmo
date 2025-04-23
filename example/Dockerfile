# 使用官方的 Golang 镜像作为基础镜像
FROM registry.cn-shanghai.aliyuncs.com/vue-gin-devops/golang1.21.1:latest

# 设置环境变量
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./
# 下载依赖
RUN go mod download
# 将项目的所有源代码复制到工作目录
COPY . .
# 空的 RUN 指令以防止缓存
RUN echo "Rebuild required"

RUN apt-get update && apt-get install -y ca-certificates

# 构建应用程序
RUN go build -o gosmo-example


# 暴露应用程序的端口
EXPOSE 9999

# 运行应用程序
CMD ["./gosmo-example"]