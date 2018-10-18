#!/bin/bash
consul agent -retry-join=consulserver &
shift
./kadfs $@