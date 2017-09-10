build:
	GOOS=linux GOARCH=amd64 go build -o bin/linux/app
	GOOS=windows GOARCH=amd64 go build -o bin/windows/app.exe
	GOOS=darwin GOARCH=amd64 go build -o bin/mac/app