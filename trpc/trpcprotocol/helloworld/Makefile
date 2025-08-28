.DEFAULT_GOAL := rpconly

rpconly:
	trpc create \
		-p helloworld.proto \
		--rpconly \
		--nogomod \
		--mock=true \
		--validate=true \

apidocs:
	trpc apidocs -p helloworld.proto

.PHONY: rpconly apidocs
