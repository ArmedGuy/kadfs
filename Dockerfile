FROM golang:latest

COPY . /go/src/github.com/ArmedGuy/kadfs
WORKDIR /go/src/github.com/ArmedGuy/kadfs

RUN go get github.com/golang/protobuf/proto
RUN go get github.com/gorilla/mux

RUN go build
RUN ./kadfs 

