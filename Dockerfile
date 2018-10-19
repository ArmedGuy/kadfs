FROM golang:latest

RUN apt-get update && apt-get install unzip
RUN wget -O /tmp/consul.zip https://releases.hashicorp.com/consul/1.3.0/consul_1.3.0_linux_amd64.zip
RUN unzip /tmp/consul.zip -d /tmp/
RUN cp /tmp/consul /bin/consul
RUN chmod +x /bin/consul
RUN go get github.com/golang/protobuf/proto
RUN go get github.com/gorilla/mux
RUN go get github.com/hashicorp/consul

COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
COPY . /go/src/github.com/ArmedGuy/kadfs
WORKDIR /go/src/github.com/ArmedGuy/kadfs




EXPOSE 4000

RUN go build

ENTRYPOINT ["/docker-entrypoint.sh"]

