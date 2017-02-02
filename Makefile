OS ?= linux
ARCH ?= amd64

build:
	GOOS=linux GOARCH=amd64 go build -o bin/dht-hallo service/cmd/dht-server/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/dht-hallo.exe service/cmd/dht-server/main.go

run:
	go run service/cmd/dht-server/main.go --port 8080

#npm install -g bootprint
#npm install -g bootprint-swagger
generate:
	#delete all files apart from the configure file.
	cd service/ && find . ! -name \configure_*.go -type f -exec rm -f {} +
	#delete empty directories
	cd service/ && find . -type d -empty -delete
	cd service/ && swagger generate server -A dht -f ../swagger.yml
	rm -rf doc/*
	bootprint swagger swagger.yml doc