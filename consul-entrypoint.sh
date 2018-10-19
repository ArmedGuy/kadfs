#!/bin/bash
IPMASK=$(ip a | grep "inet 10" | awk '{print $2}')
consul agent -server -bind="${IPMASK::-3}" $@
