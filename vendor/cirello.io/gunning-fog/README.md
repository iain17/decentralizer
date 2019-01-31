#gunning-fog

Gunning-fog index analyzer written in Go. This analyzer processes an English
text and produces its Gunning Fox index score. Refer to its logic in
https://en.wikipedia.org/wiki/Gunning_fog_index - it does not analyse word
endings (-es, -ed, or -ing), or discriminate proper nouns, familiar jargon or
compound words.

*Use*

```ShellSession
$ go get cirello.io/gunning-fog/...
$ cat LICENSE | $GOPATH/bin/gunning-fog
16
```

`gunning-fog` will always wait content from STDIN.