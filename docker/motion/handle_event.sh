#!/usr/bin/env bash

set -e -x

echo "[$(date --iso-8601=ns), ${*}]" | /bin/nc -q 0 -u "${UDP_HOST}" "${UDP_PORT}"
