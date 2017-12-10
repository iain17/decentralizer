TARGET=adna
build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o bin/linux/$(TARGET) main.go
	upx --brute bin/linux/$(TARGET)
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o bin/windows/$(TARGET).exe main.go
	upx --brute bin/windows/$(TARGET).exe
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o bin/mac/$(TARGET) main.go
	upx --brute bin/mac/$(TARGET)

#https://github.com/jteeuwen/go-bindata
generate:
	./generate.sh