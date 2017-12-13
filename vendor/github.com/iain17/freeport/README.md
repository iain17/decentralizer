FreePort
========

Get a free open TCP or UDP port that is ready to use

```bash
# Ask the kernel to give us an open port.
export port=$(freeport TCP)

# Start standalone httpd server for testing
httpd -X -c "Listen $port" &

# Curl local server on the selected port
curl localhost:$port
```

#### usage
```bash
sudo apt-get install golang                    # Download go. Alternativly build from source: https://golang.org/doc/install/source
mkdir ~/.gopath && export GOPATH=~/.gopath     # Replace with desired GOPATH
export PATH=$PATH:$GOPATH/bin                  # For convenience, add go's bin dir to your PATH
go get github.com/iain17/freeport/cmd/freeport
```
