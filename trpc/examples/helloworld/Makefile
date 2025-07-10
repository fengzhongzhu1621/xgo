rpconly:
	trpc create \
		-p helloworld.proto \
		--rpconly \
		--nogomod \
		--mock=false \

all:
	trpc create -p helloworld.proto


build:
	go build

.PHONY: rpconly all
