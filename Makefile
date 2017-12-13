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

test:
	go test $(go list ./... | grep -v /vendor/)

#https://github.com/jteeuwen/go-bindata
generate:
	./generate.sh