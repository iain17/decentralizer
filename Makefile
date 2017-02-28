build:
	GOOS=linux GOARCH=amd64 go build -o bin/decentralizer service/cmd/dht-server/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/decentralizer.exe service/cmd/dht-server/main.go

run:
	go run service/cmd/decentralizer-server/main.go --port 8080

#npm install -g bootprint
#npm install -g bootprint-swagger
generate:
	#delete all files apart from the configure file.
	cd service/ && find . ! -name \configure_*.go -type f -exec rm -f {} +
	#delete empty directories
	cd service/ && find . -type d -empty -delete
	cd service/ && swagger generate server -A decentralizer -f ../swagger.yml
	rm -rf doc/*
	bootprint swagger swagger.yml doc