#!/bin/bash

set -e

CF_LOGIN_SCRIPT=$(hammer cf-login --file | tail -n1)
# shellcheck source=/dev/null
. "${CF_LOGIN_SCRIPT}" >/dev/null

CF_API_URL=$(cf api | head -n1 | awk '{print $NF}')

DEPLOYMENT=$(hammer bosh -- deployments --json | jq -r .Tables[0].Rows[0].name)
CF_WORKER_IP=$(hammer bosh -- --deployment "${DEPLOYMENT}" ssh compute --command 'curl\ checkip.amazonaws.com' | grep stdout | awk '{print $NF}' | tr -dc '0-9.')

cat >./test/e2e/tanzu_data <<EOL
export CF_WORKER_IP=${CF_WORKER_IP}
export CF_USERNAME=${CF_USERNAME}
export CF_PASSWORD=${CF_PASSWORD}
export CF_API_URL=${CF_API_URL}
EOL
