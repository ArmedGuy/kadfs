#!/bin/bash
IPMASK=$(ip a | grep "inet 10" | awk '{print $2}')
consul agent -server -advertise="${IPMASK::-3}" $@ &
FABIO_LOG_ACCESS_TARGET=stdout FABIO_PROXY_DIALTIMEOUT=300s FABIO_LOG_ACCESS_FORMAT='$remote_host - $upstream_host [$time_common] "$request" $response_status $response_body_size' fabio
