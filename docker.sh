#!/bin/bash

MODULE_NAME=xgo
CONTAINER_NAME=xgo
VERSION=$(date +%Y%m%d%H%M)

CLOUD_DOMAIN=mirrors.tencent.com
DOCKER_NAMESPACE=xgo

rm_image() {
    if [ -z $1 ]; then
        echo "docker rmi failed, module name not defined"
        exit 1
    fi

    docker images | grep $1
    if [ $? = 0 ]; then
        echo "docker rmi $1"
        docker images | grep $1 | awk '{print $1":"$2}' | xargs docker rmi
    fi
}

rm_container() {
    if [ -z $1 ]; then
        echo "docker rm failed, container name not defined"
        exit 1
    fi

    echo "docker stop $1"
    docker stop $1

    echo "docker rm $1"
    docker rm $1
}

docker_build() {
    rm_container ${CONTAINER_NAME}
    rm_image ${MODULE_NAME}

    echo "docker build -t ${MODULE_NAME}:${VERSION} ."
    docker build -t ${MODULE_NAME}:${VERSION} .
}

docker_push() {
    echo "docker login"
    docker login ${CLOUD_DOMAIN} || exit 1

    echo "docker tag ${name} ${CLOUD_DOMAIN}/${DOCKER_NAMESPACE}/${name}"
    name=$(docker images | grep $MODULE_NAME | awk '{print $1":"$2}')
    docker tag ${name} ${CLOUD_DOMAIN}/${DOCKER_NAMESPACE}/${name}

    echo "docker push ${CLOUD_DOMAIN}/${DOCKER_NAMESPACE}/${name}"
    docker push ${CLOUD_DOMAIN}/${DOCKER_NAMESPACE}/${name}

    echo "docker rmi ${CLOUD_DOMAIN}/${DOCKER_NAMESPACE}/${name}"
    docker rmi ${CLOUD_DOMAIN}/${DOCKER_NAMESPACE}/${name}
}

main() {
    docker_build
    docker_push
}

main
