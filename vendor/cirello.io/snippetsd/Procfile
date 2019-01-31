workdir: $GOPATH/src/cirello.io/snippetsd
observe: *.go
ignore: /vendor
build-backed: go install cirello.io/snippetsd/cmd/snippetsd
backend:      $GOPATH/bin/snippetsd http
ui:           restart=tmp cd frontend; npm run start

