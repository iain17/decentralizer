workdir: $GOPATH/src/cirello.io/bookmarkd
observe: *.go
ignore: /vendor
build-backend: CC=gcc vgo install cirello.io/bookmarkd/cmd/bookmarkd
backend:       $GOPATH/bin/bookmarkd http
ui:            restart=tmp cd frontend; npm run start
