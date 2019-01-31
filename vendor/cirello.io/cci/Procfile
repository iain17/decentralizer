workdir: $GOPATH/src/cirello.io/cci
observe: *.go *.yaml
ignore: /vendor
build-backend: CC=gcc vgo install cirello.io/cci/cmd/cci
server:        $GOPATH/bin/cci standalone
