#!/bin/sh

docker build . -t kademlia
docker-compose -f bootstrap-compose.yml up -d

for i in `seq 1 5`;
do
    docker-compose -f kademlia-compose.yml up -d
done
