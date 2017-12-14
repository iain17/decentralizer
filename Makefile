TARGET=adna
#docker run -d -v /Users/iain17/work/src/github.com/iain17/decentralizer/:/app -i golang

#apt-get -y update
#apt-get -y install build-essential upx-ucl
#go get -v -u github.com/whyrusleeping/gx
#go get -v -u github.com/golang/dep/cmd/dep

install:
	$(GOPATH)/bin/gx install
	ls $(GOPATH)/src/gx/ipfs
    #Patch a stupid fucking problem because of gx and the way ipfs does deps: debug/requests problem
	find $(GOPATH)/src/gx/ -name 'trace.go' -exec sed -i '.bak' -e 's/requests"/requestss"/g' {} \;
	find $(GOPATH)/src/gx/ -name 'trace.go' -exec sed -i '.bak' -e 's/events"/eventss"/g' {} \;
	find $(GOPATH)/src/gx/ -name '*.bak' -type f -exec rm -f {} +
	$(GOPATH)/bin/dep ensure

build-linux:
	GOOS=linux GOARCH=amd64 go build -gcflags=-trimpath=x/y -ldflags "-s -w" -o bin/linux/$(TARGET) main.go
	cp bin/linux/$(TARGET) bin/linux/unpacked-$(TARGET)
	#upx --brute bin/linux/$(TARGET)

build-win:
	GOOS=windows GOARCH=amd64 go build -gcflags=-trimpath=x/y -ldflags "-s -w" -o bin/windows/$(TARGET).exe main.go
	cp bin/windows/$(TARGET).exe bin/windows/unpacked-$(TARGET).exe
	#upx --brute bin/windows/$(TARGET).exe

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -gcflags=-trimpath="github.com/iain17" -ldflags "-s -w" -o bin/mac/$(TARGET) main.go
	cp bin/mac/$(TARGET) bin/mac/unpacked-$(TARGET)
	#upx --brute bin/mac/$(TARGET)

ci:
	gitlab-runner --debug exec docker test

test:
	go test -race -cover ./...

#https://github.com/jteeuwen/go-bindata
generate:
	./generate.sh