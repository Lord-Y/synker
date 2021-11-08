#!/bin/bash

set -x

/cockroach/cockroach sql --insecure --host ${COCKROACH_HOST} --execute "SELECT 1" || /cockroach/cockroach init --insecure --host ${COCKROACH_HOST}

exit 0
