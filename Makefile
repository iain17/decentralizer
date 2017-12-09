TARGET=chat

build:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/$(TARGET) main.go
	#GOOS=windows GOARCH=amd64 go build -o bin/windows/$(TARGET).exe main.go
	#GOOS=darwin GOARCH=amd64 go build -o bin/mac/$(TARGET) main.go

#https://github.com/jteeuwen/go-bindata
generate:
	protoc --go_out=. discovery/pb/protocol.proto
	protoc --go_out=. pb/*.proto
	go-bindata -pkg static -o static/data.go static/