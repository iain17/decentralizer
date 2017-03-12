build:
	GOOS=linux GOARCH=amd64 go build -o bin/decentralizer main.go
	GOOS=windows GOARCH=amd64 go build -o bin/decentralizer.exe main.go

run:
	go run main.go serve --listen :8080