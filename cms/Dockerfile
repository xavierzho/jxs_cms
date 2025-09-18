# ---------- 构建阶段 ----------
FROM golang:1.24.5 AS  builder

# 安装构建需要的工具

WORKDIR /app

# 先复制 go.mod 和 go.sum，利用缓存
COPY go.mod go.sum ./

# 再复制源码
COPY . .

# 设置环境变量，关闭 CGO，编译 Linux amd64 二进制
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOPROXY=https://goproxy.cn go build -o chaoshe_data_backend .

# ---------- 运行阶段 ----------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache tzdata \
 && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
 && echo "Asia/Shanghai" > /etc/timezone

# 从 builder 阶段复制编译好的二进制
COPY --from=builder /app/chaoshe_data_backend ./data_backend

# 复制运行需要的配置和资源
COPY ./configs ./configs
COPY ./storage ./storage
COPY ./start.sh ./start.sh

RUN chmod +x ./start.sh ./data_backend

EXPOSE 8003

CMD ["/app/start.sh"]
