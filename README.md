# go-exercises-worker

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-exercises-worker?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-exercises-worker)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-exercises-worker)](https://goreportcard.com/report/github.com/thewizardplusplus/go-exercises-worker)
[![Build Status](https://travis-ci.com/thewizardplusplus/go-exercises-worker.svg?branch=master)](https://travis-ci.com/thewizardplusplus/go-exercises-worker)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-exercises-worker/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-exercises-worker)

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

## Usage

```
$ go-exercises-backend
```

Environment variables:

- `ALLOWED_IMPORT_CONFIG` &mdash; path to the allowed import config (default: `./configs/allowed_imports.json`);
- `RUNNING_TIMEOUT` &mdash; maximal duration of solution running (default: `10s`);
- message broker:
  - `MESSAGE_BROKER_ADDRESS` &mdash; [RabbitMQ](https://www.rabbitmq.com/) connection URI (default: `amqp://rabbitmq:rabbitmq@localhost:5672`);
  - `MESSAGE_BROKER_BUFFER_SIZE` &mdash; [RabbitMQ](https://www.rabbitmq.com/) channel capacity (default: `1000`);
- `SOLUTION_CONSUMER_CONCURRENCY` &mdash; amount of solution consumer threads (default: `1000`).

## API Description

API description in the [AsyncAPI](https://www.asyncapi.com/) format: [docs/async_api.yaml](docs/async_api.yaml).

## License

The MIT License (MIT)

Copyright &copy; 2021 thewizardplusplus
