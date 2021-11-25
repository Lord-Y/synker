#!/bin/bash

set -euo pipefail

COCKROACH_HOST=$(netstat -latn |grep 26257 |grep LISTEN |awk '{print $4}')
cockroach sql --insecure --host ${COCKROACH_HOST} --execute "SHOW CLUSTER SETTING kv.rangefeed.enabled" --format tsv | \
  grep true || \
  cockroach sql --insecure --host ${COCKROACH_HOST} --execute "SET CLUSTER SETTING kv.rangefeed.enabled = true"


cockroach sql --insecure --host ${COCKROACH_HOST} --execute "SHOW CHANGEFEED JOBS" --format tsv | \
  grep kafka || \
  cockroach sql --insecure --host ${COCKROACH_HOST} --execute "CREATE CHANGEFEED FOR TABLE promo_codes, rides, user_promo_codes, users, vehicle_location_histories, vehicles INTO 'kafka://localhost:9092' WITH full_table_name, on_error = 'pause', protect_data_from_gc_on_pause, updated, resolved" --database movr
