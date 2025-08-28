build:
	go mod tidy
	go build -o trpc.examples.helloworld.Greeter

run: build
	./trpc.examples.helloworld.Greeter

docker:
	docker build -t trpc.examples.helloworld.greeter:latest .

generate:
	go generate ./...
