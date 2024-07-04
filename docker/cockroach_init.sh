#!/bin/bash

set -euo pipefail

COCKROACH_HOST=$(netstat -latn |grep 26257 |grep LISTEN |awk '{print $4}')
cockroach sql --insecure --host ${COCKROACH_HOST} --execute "SHOW CLUSTER SETTING kv.rangefeed.enabled" --format tsv | \
  grep true || \
  cockroach sql --insecure --host ${COCKROACH_HOST} --execute "SET CLUSTER SETTING kv.rangefeed.enabled = true"
