FROM golang:1.15-alpine AS builder

RUN apk update && \
  apk add --no-cache curl git && \
  curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.4/dep-linux-amd64 && \
  chmod +x /usr/local/bin/dep && \
  go get golang.org/x/tools/cmd/goimports

WORKDIR /go/src/github.com/thewizardplusplus/go-exercises-worker
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure -vendor-only -v

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -a -ldflags='-w -s -extldflags "-static"' ./...

FROM golang:1.15-alpine

RUN apk update && apk add --no-cache bash

COPY --from=builder /go/bin/go-exercises-worker /go/bin/goimports /usr/local/bin/
COPY tools/wait-for-it.sh /usr/local/bin/wait-for-it.sh

CMD ["/usr/local/bin/go-exercises-worker"]
