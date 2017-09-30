# go-stubfilegenerator utility

[![license](https://img.shields.io/github/license/webdevops/go-stubfilegenerator.svg)](https://github.com/webdevops/go-stubfilegenerator/blob/master/LICENSE)

Stub file generator written in golang

Features
========
- Creates text files (txt, csv xml, css, js..)
- Creates images (jpg, gif, png)
- Creates documents (pdf, rtf, doc)
- Creates archives (zip, tar, tar.bz2, tar.gz)

Usage
=====

Install: ```go get -u github.com/webdevops/go-stubfilegenerator```

```go
package main

import (
	"github.com/webdevops/go-stubfilegenerator"
)

func main() {
	stubgen := stubfilegenerator.StubGenerator()
	stubgen.GenerateStub("example/foo.txt")
	stubgen.GenerateStub("example/foo.csv")
	stubgen.GenerateStub("example/foo.xml")
	stubgen.GenerateStub("example/foo.gif")
	stubgen.GenerateStub("example/foo.png")
	stubgen.GenerateStub("example/foo.jpg")

	stubgen.Image.Width = 1000
	stubgen.Image.Height = 1000
	stubgen.GenerateStub("example/big.jpg")
}
```
