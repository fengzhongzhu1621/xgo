# -------------- builder container --------------
# 用于编译 go 二进制文件
FROM golang:1.23-alpine as builder

WORKDIR /app

ARG VERSION

COPY go.mod .
COPY go.sum .

COPY . .

RUN go env -w GO111MODULE=on
# RUN go env -w GOPROXY=https://mirrors.cloud.tencent.com/go/,direct \
RUN go env -w GOPROXY=https://goproxy.cn,direct
run go mod download



# 编译
RUN  go build -o ./xgo main/main.go

# -------------- runner container --------------
# Alpine 操作系统是一个面向安全的轻型 Linux 发行版。它不同于通常 Linux 发行版，
# Alpine 采用了 musl libc 和 busybox 以减小系统的体积和运行时资源消耗，但功能上比 busybox 又完善的多，因此得到开源社区越来越多的青睐。
# 在保持瘦身的同时，Alpine 还提供了自己的包管理工具 apk，可以通过 https://pkgs.alpinelinux.org/packages 网站上查询包信息，
# 也可以直接通过 apk 命令直接查询和安装各种软件。
#
# Alpine Docker 镜像也继承了 Alpine Linux 发行版的这些优势。相比于其他 Docker 镜像，它的容量非常小，仅仅只有 5 MB 左右
# （对比 Ubuntu 系列镜像接近 200 MB），且拥有非常友好的包管理机制。官方镜像来自 docker-alpine 项目。
#
# 目前 Docker 官方已开始推荐使用 Alpine 替代之前的 Ubuntu 做为基础镜像环境。
# 这样会带来多个好处。包括镜像下载速度加快，镜像安全性提高，主机之间的切换更方便，占用更少磁盘空间等。
FROM alpine:latest AS runner

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tencent.com/g' /etc/apk/repositories

RUN apk --update --no-cache add bash
RUN apk --no-cache add tzdata

# 设置工作目录
WORKDIR /app

# RUN 指令用于执行一个或多个命令。如果需要执行多个命令，可以将它们连接起来，使用 && 符号分隔。
RUN mkdir -p /app/docs \
    && mkdir -p /app/logs \
    && mkdir -p /app/templates \
    && mkdir -p /app/static \
    && mkdir -p /app/logs \
    && mkdir -p app/config

# 设置多个环境变量
ENV TMPLATE_FILE_BASE_DIR=/app/templates \
    DOC_FILE_BASE_DIR=/app/docs \
    STATIC_FILE_BASE_DIR=/app/static \
    LOG_FILE_BASE_DIR=/app/logs

# 复制应用程序文件到工作目录
COPY --from=builder /app/xgo /usr/bin/xgo
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/static /app/static
COPY --from=builder /app/config/config.yaml /app/config/

EXPOSE 8080
CMD ["/app/xgo"]
