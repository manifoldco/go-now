# go-now

Client for [Zeit Now](https://zeit.co/api) deployment API

[Code of Conduct](./CODE_OF_CONDUCT.md) |
[Contribution Guidelines](./.github/CONTRIBUTING.md)

[![GitHub release](https://img.shields.io/github/tag/manifoldco/go-now.svg?label=latest)](https://github.com/manifoldco/go-now/releases)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/manifoldco/go-now)
[![Travis](https://img.shields.io/travis/manifoldco/go-now/master.svg)](https://travis-ci.org/manifoldco/go-now)
[![Go Report Card](https://goreportcard.com/badge/github.com/manifoldco/go-now)](https://goreportcard.com/report/github.com/manifoldco/go-now)
[![License](https://img.shields.io/badge/license-BSD-blue.svg)](./LICENSE.md)

## Installation

```
go get github.com/manifoldco/go-now
```

## Example

```go
import "github.com/manifoldco/go-now"

n := now.New("your-api-secret")

pkg := map[string]interface{}{
  "index.js": "require('http').Server((req, res) => { res.end('Hello World!'); }).listen();",
  "package": map[string]interface{}{
    "name": "hello-world",
    "scripts": map[string]string{
      "start": "node index",
    },
  },
}
d, err := n.Deployment.New(pkg)

// &{UID: "7Npest0z1zW5QVFfNDBId4BW", Host: "hello-world-abcdefhi.now.sh", State: "BOOTING"} 
```
