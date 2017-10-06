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
	stubgen := stubfilegenerator.NewStubGenerator()
	stubgen.Generate("example/foo.txt")
	stubgen.Generate("example/foo.csv")
	stubgen.Generate("example/foo.xml")
	stubgen.Generate("example/foo.gif")
	stubgen.Generate("example/foo.png")
	stubgen.Generate("example/foo.jpg")

	stubgen.Image.Width = 1000
	stubgen.Image.Height = 1000
	stubgen.Generate("example/big.jpg")
}
```
