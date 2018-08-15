# Update IPFS
Firstly know that this is going to take quite long. But here are some hints on how to get it done.
1. Download a IPFS [release](https://github.com/ipfs/go-ipfs/releases).
2. cd into the release. And run: ```make install```
3. ```gx publish``` and get the has for the base IPFS repo.
4. Copy all the files from: ```$GOPATH/src/gx``` into this projects vendor/gx directory.
5. Update all the import paths.