TARGET=adna
install:
	./install.sh

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/linux/$(TARGET) main.go
	upx --brute bin/linux/$(TARGET)

build-win:
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/windows/$(TARGET).exe main.go
	upx --brute bin/windows/$(TARGET).exe

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/mac/$(TARGET) main.go
	upx --brute bin/mac/$(TARGET)

ci:
	gitlab-runner --debug exec docker test

#https://github.com/jteeuwen/go-bindata
generate:
	./generate.sh