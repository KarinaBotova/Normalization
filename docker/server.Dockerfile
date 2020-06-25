FROM golang:alpine

RUN apk add --no-cache make

WORKDIR /go/src/github.com/kolya59/easy-normalization

COPY . .

RUN make build-server-linux && mv ./bin/linux/server.bin /exe

WORKDIR /
ENTRYPOINT ["/exe"]