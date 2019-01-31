all:
	protoc -I$(GOPATH)/src:. --gofast_out=plugins=grpc:. ./*.proto
	gomodifytags -w -file api.pb.go -struct Recipe -add-tags yaml,db 1> /dev/null
	gomodifytags -w -file api.pb.go -struct Build -add-tags yaml,db 1> /dev/null