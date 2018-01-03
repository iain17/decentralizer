TARGET:=adna
ARCH:=windows#darwin, linux
#docker run -d -v /Users/iain17/work/src/github.com/iain17/decentralizer/:/app -i golang

#apt-get -y update
#apt-get -y install build-essential upx-ucl
#go get -v -u github.com/whyrusleeping/gx
#go get -v -u github.com/golang/dep/cmd/dep

install:
	$(GOPATH)/bin/gx install
	ls $(GOPATH)/src/gx/ipfs
    #Patch a stupid fucking problem because of gx and the way ipfs does deps: debug/requests problem
	find vendor/gx/ -name 'trace.go' -exec sed -i '.bak' -e 's/requests"/requestss"/g' {} \;
	find vendor/gx/ -name 'trace.go' -exec sed -i '.bak' -e 's/events"/eventss"/g' {} \;
	find vendor/gx/ -name '*.bak' -type f -exec rm -f {} +
	$(GOPATH)/bin/dep ensure

build:
	mkdir -p bin/$(ARCH)/
	rm -rf bin/$(ARCH)/*
	GOOS=$(ARCH) GOARCH=amd64 go build -ldflags "-s -w" -o bin/$(ARCH)/unpacked-$(TARGET) main.go

pack:
	rm bin/$(ARCH)/$(TARGET)
	upx --brute -o bin/$(ARCH)/$(TARGET) bin/$(ARCH)/unpacked-$(TARGET)

ci:
	gitlab-runner --debug exec docker test

test:
	go test -cover ./...
gx:
	rm -rf vendor/gx
	gx install

dep:
	mv vendor/gx /tmp/
	dep ensure
	mv /tmp/gx vendor/

dep-update:
	mv vendor/gx /tmp/
	dep ensure --update
	dep prune
	mv /tmp/gx vendor/

#https://github.com/jteeuwen/go-bindata
generate:
	./generate.sh