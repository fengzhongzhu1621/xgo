# 第一阶段：构建阶段
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o trpc.examples.helloworld.Greeter .

# 第二阶段：运行阶段
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/trpc.examples.helloworld.Greeter /app/hellworld
COPY --from=builder /app/trpc_go.yaml /app/trpc_go.yaml

EXPOSE 8000
CMD ["/app/hellworld"]
