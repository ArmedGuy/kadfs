#!/bin/sh

docker build . -t kademlia
docker-compose -f bootstrap-compose.yml up -d

for i in `seq 1 5`;
do
    docker run --network=testing_kademlia_network -itd kademlia 4000 0000000000000000000000000000000000000000 172.16.238.10:4000
done
