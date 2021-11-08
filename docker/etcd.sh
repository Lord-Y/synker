#!/bin/bash

etcd --name ${NODE_NAME} \
  --initial-advertise-peer-urls https://${NODE_NAME}:2380 \
  --listen-peer-urls https://0.0.0.0:2380 \
  --listen-client-urls https://0.0.0.0:2379 \
  --advertise-client-urls https://${NODE_NAME}:2379 \
  --initial-cluster-token ${INITIAL_CLUSTER_TOKEN} \
  --initial-cluster etcdc01=https://${NODE_NAME1}:2380,etcdc02=https://${NODE_NAME2}:2380,etcdc03=https://${NODE_NAME3}:2380 \
  --initial-cluster-state ${INITIAL_CLUSTER_STATE} \
  --auto-tls \
  --peer-auto-tls \
  --data-dir ${DATA_DIR}
