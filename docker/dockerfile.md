# 指令
## LABEL

用于为 Docker 镜像添加元数据。这个指令允许你指定一个联系人，以便在出现问题时能够找到负责维护该镜像的人。

```
LABEL key=value
```

```
FROM ubuntu:latest
LABEL maintainer="John Doe <john.doe@example.com>"
```


这将使得在构建的 Docker 镜像中包含以下元数据：
```
{
  "Maintainer": "John Doe <john.doe@example.com>"
}
```


## COPY

用于将文件或目录从构建上下文复制到 Docker 镜像中，以便在容器中使用。

```
COPY . /app/gohttpserver/
COPY a.tar.gz /tmp
COPY . /
COPY src /app/src
COPY app.py /app/
```


## ENV
用于设置环境变量。这些环境变量可以在容器运行时被应用程序使用，或者在 Dockerfile 中的其他指令中使用。

```
ENV key value
ENV APP_HOME /app
ENV TIME_ZONE=Asia/Shanghai
ENV DB_HOST=localhost \
    DB_PORT=5432 \
    DB_USER=myuser \
    DB_PASS=mypassword
```

在 Dockerfile 中使用环境变量
```
# 使用 ENV 指令设置的环境变量
COPY $APP_HOME/app.py /app/

# 在 RUN 指令中使用环境变量
RUN echo "Database host: $DB_HOST"
```


在容器运行时使用环境变量

当你使用 docker run 命令运行容器时，可以通过 -e 选项传递环境变量：
```
docker run -e DB_HOST=myhost myimage
```
这将覆盖 Dockerfile 中设置的 DB_HOST 环境变量的值。


## RUN
