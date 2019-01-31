test:
	go get -u golang.org/x/vgo
	vgo test -v ./pkg/...

linux:
	docker run -ti --rm -v $(GOPATH):/go/ \
		-e CC=gcc \
		-w /go/src/cirello.io/cci golang \
		/bin/bash -c 'go get -u golang.org/x/vgo && vgo build -o cci.linux ./cmd/cci'

