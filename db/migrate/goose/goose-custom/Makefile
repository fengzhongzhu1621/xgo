build:
	go build -o goose-custom
	chmod +x goose-custom

run: build
	./goose-custom -dir . mysql "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local" up

status:
	./goose-custom -dir . mysql "root:@tcp(127.0.0.1:3306)/xgo?parseTime=True&loc=Local" status
