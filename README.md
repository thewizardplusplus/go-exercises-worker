# go-exercises-worker

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-exercises-worker?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-exercises-worker)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-exercises-worker)](https://goreportcard.com/report/github.com/thewizardplusplus/go-exercises-worker)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-exercises-worker.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-exercises-worker)

## Installation

Prepare the directory:

```
$ mkdir --parents "$(go env GOPATH)/src/github.com/thewizardplusplus/"
$ cd "$(go env GOPATH)/src/github.com/thewizardplusplus/"
```

Clone this repository:

```
$ git clone https://github.com/thewizardplusplus/go-exercises-worker.git
$ cd go-exercises-worker
```

Install dependencies with the [dep](https://golang.github.io/dep/) tool:

```
$ dep ensure -vendor-only
```

Build the project:

```
$ go install ./...
```

## License

The MIT License (MIT)

Copyright &copy; 2021 thewizardplusplus
