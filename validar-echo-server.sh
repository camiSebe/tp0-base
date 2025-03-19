#!/bin/bash

NETWORK="tp0_testing_net"
TTL=2
SERVER_CONTAINER_NAME="server"
SERVER_PORT="12345"
DUMMY_MESSAGE="ping"


RESPUESTA=$(docker run --rm --network $NETWORK busybox nc -w $TTL $SERVER_CONTAINER_NAME $SERVER_PORT <<< "$DUMMY_MESSAGE")

if [[ "$RESPUESTA" == "$DUMMY_MESSAGE" ]]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
