# kvcache

[![Build Status](https://travis-ci.org/Akagi201/kvcache.svg)](https://travis-ci.org/Akagi201/kvcache) [![Coverage Status](https://coveralls.io/repos/github/Akagi201/kvcache/badge.svg?branch=master)](https://coveralls.io/github/Akagi201/kvcache?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/Akagi201/kvcache)](https://goreportcard.com/report/github.com/Akagi201/kvcache) [![GoDoc](https://godoc.org/github.com/Akagi201/kvcache?status.svg)](https://godoc.org/github.com/Akagi201/kvcache)

A distributed in-memory key:val cache

## Import

* `import "github.com/Akagi201/kvcache/ttlru"` full version, implemented with goroutine.
* `import "github.com/Akagi201/kvcache/lttlru"` light version, implemented without goroutine.
