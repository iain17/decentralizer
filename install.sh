#!/bin/bash
$(GOPATH)/bin/gx install
#Patch a stupid fucking problem because of gx and the way ipfs does deps: debug/requests problem
find $(GOPATH)/src/gx/ -name 'trace.go' -exec sed -i '.bak' -e 's/requests"/requestss"/g' {} \;
find $(GOPATH)/src/gx/ -name 'trace.go' -exec sed -i '.bak' -e 's/events"/eventss"/g' {} \;
find $(GOPATH)/src/gx/ -name '*.bak' -type f -exec rm -f {} +
$(GOPATH)/bin/dep ensure