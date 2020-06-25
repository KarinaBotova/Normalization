FROM golang:alpine

RUN make build-client-linux
COPY ./bin/linux/client.bin /exe

WORKDIR /
ENTRYPOINT ["/exe"]