FROM registry.digitalarsenal.net/docker-images/general-go:1.10.0

ENV GOPATH /go
ENV PATH $PATH:/go/bin

RUN apt update
RUN apt-get install -y protobuf-compiler

WORKDIR /app
