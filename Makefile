TARGET:=adna
ARCH:=amd64
GOOSE:=windows#darwin, linux
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

clean:
	rm -rf bin/$(GOOSE)/

build:
	mkdir -p bin/$(GOOSE)/$(ARCH)
	GOOS=$(GOOSE) GOARCH=$(ARCH) go build -ldflags "-s -w" -o bin/$(GOOSE)/$(ARCH)/unpacked-$(TARGET) main.go

pack:
	rm bin/$(ARCH)/$(TARGET) || true
	upx -o bin/$(GOOSE)/$(ARCH)/$(TARGET) bin/$(GOOSE)/$(ARCH)/unpacked-$(TARGET)
	rm bin/$(GOOSE)/$(ARCH)/unpacked-$(TARGET) || true

ci:
	gitlab-runner --debug exec docker build

test:
	go test -v -cover ./...
gx:
	rm -rf vendor/gx
	gx install

dep:
	mv vendor/gx /tmp/
	dep ensure
	mv /tmp/gx vendor/

dep-update:
	mv vendor/gx /tmp/
	dep ensure
	dep ensure --update
	mv /tmp/gx vendor/

#https://github.com/jteeuwen/go-bindata
generate:
	./generate.sh