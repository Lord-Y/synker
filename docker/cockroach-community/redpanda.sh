#!/bin/bash

set -x

if [ -n "${LEADER}" ]
then
    rpk redpanda start \
    --smp 1  \
    --memory 1G  \
    --reserve-memory 0M \
    --overprovisioned \
    --node-id ${NODE_ID} \
    --check=false \
    --pandaproxy-addr 0.0.0.0:8082 \
    --advertise-pandaproxy-addr 127.0.0.1:8082 \
    --kafka-addr 0.0.0.0:9092 \
    --advertise-kafka-addr 127.0.0.1:9092 \
    --rpc-addr 0.0.0.0:33145 \
    --advertise-rpc-addr ${NODE_NAME}:33145
else
    rpk redpanda start \
    --smp 1  \
    --memory 1G  \
    --reserve-memory 0M \
    --overprovisioned \
    --node-id ${NODE_ID} \
    --seeds "${SEED_NODE_NAME}:33145" \
    --check=false \
    --pandaproxy-addr 0.0.0.0:8082 \
    --advertise-pandaproxy-addr 127.0.0.1:8082 \
    --kafka-addr 0.0.0.0:9092 \
    --advertise-kafka-addr 127.0.0.1:9092 \
    --rpc-addr 0.0.0.0:33145 \
    --advertise-rpc-addr ${NODE_NAME}:33145
fi
