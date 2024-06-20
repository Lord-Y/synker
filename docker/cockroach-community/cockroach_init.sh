#!/bin/bash

set -euo pipefail

# https://www.cockroachlabs.com/docs/v24.1/start-a-local-cluster-in-docker-linux

sudo docker exec -ti cockroachc01 ./cockroach sql --insecure --host localhost:26257 --execute "SELECT 1" || sudo docker exec -ti cockroachc01 ./cockroach init --insecure --host localhost:26357

COCKROACH_HOST=$(netstat -latn |grep 26257 |grep LISTEN |awk '{print $4}' | head -1)
cockroach sql --insecure --host ${COCKROACH_HOST} --execute "SHOW CLUSTER SETTING kv.rangefeed.enabled" --format json | jq .[].'"kv.rangefeed.enabled"' | \
  grep '"t"' || \
  cockroach sql --insecure --host ${COCKROACH_HOST} --execute "SET CLUSTER SETTING kv.rangefeed.enabled = true"
