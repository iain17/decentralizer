all: assets darwin

assets:
	cd frontend; npm install; npm run build;
	go-bindata-assetfs -o bindata_assetfs.go -pkg generated frontend/build/...
	mv bindata_assetfs.go generated

darwin:
	vgo build -o bookmarkd ./cmd/bookmarkd

linux:
	docker run -ti --rm -v $(PWD)/../:/go/src/cirello.io/ \
		-w /go/src/cirello.io/bookmarkd golang \
		/bin/bash -c 'go get -u golang.org/x/vgo && vgo build -o bookmarkd.linux ./cmd/bookmarkd'

test:
	go get -u golang.org/x/vgo
	vgo test -v ./...
