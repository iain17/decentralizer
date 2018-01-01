TARGET:=discovery
ARCH:=windows#darwin, linux

build:
	GOOS=$(ARCH) GOARCH=amd64 go build -ldflags "-s -w" -o bin/$(TARGET) main.go