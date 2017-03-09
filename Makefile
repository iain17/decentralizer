build:
	GOOS=linux GOARCH=amd64 go build -o bin/decentralizer service/cmd/dht-server/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/decentralizer.exe service/cmd/dht-server/main.go

run:
	go run main.go serve --listen :8080