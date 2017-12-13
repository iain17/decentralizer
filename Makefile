TARGET=adna
install:
	$(GOPATH)/bin/gx install
	ls $(GOPATH)/src/gx/ipfs
    #Patch a stupid fucking problem because of gx and the way ipfs does deps: debug/requests problem
	find ./vendor/gx/ -name 'trace.go' -exec sed -i '.bak' -e 's/requests"/requestss"/g' {} \;
	find ./vendor/gx/ -name 'trace.go' -exec sed -i '.bak' -e 's/events"/eventss"/g' {} \;
	find ./vendor/gx/ -name '*.bak' -type f -exec rm -f {} +
	$(GOPATH)/bin/dep ensure
	go get ./...

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