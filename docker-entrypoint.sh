#!/bin/bash
IPMASK=$(ip a | grep "inet 10" | awk '{print $2}')     
consul agent -bind="${IPMASK::-3}" -retry-join=consulserver -data-dir=/tmp &
shift
./kadfs $@
