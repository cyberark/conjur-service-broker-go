#!/bin/bash

set -e
cd "$(dirname "$0")"

# allow overwriting the hammerfile.json location
export HAMMER_TARGET_CONFIG="${HAMMER_TARGET_CONFIG:-"$PWD/../../hammerfile.json"}"

CF_LOGIN_SCRIPT=$(hammer cf-login --file | tail -n1)
# shellcheck source=/dev/null
. "${CF_LOGIN_SCRIPT}" >/dev/null

CF_API_URL=$(cf api | head -n1 | awk '{print $NF}')

# execute ssh command on the actual tanzu worker to get its IP address, to be able to allow this machine to call our aws conjur test instance
DEPLOYMENT=$(hammer bosh -- deployments --json | jq -r .Tables[0].Rows[0].name)
CF_WORKER_IP=$(hammer bosh -- --deployment "${DEPLOYMENT}" ssh compute --command 'curl\ checkip.amazonaws.com' | grep stdout | awk '{print $NF}' | tr -dc '0-9.')

if [ -z "${CF_WORKER_IP}" ]; then
  echo "failed to get worker IP address!"
  exit 1
fi

# store the tanzu env details in a file
cat >./tanzu_data <<EOL
export CF_WORKER_IP=${CF_WORKER_IP}
export CF_USERNAME=${CF_USERNAME}
export CF_PASSWORD=${CF_PASSWORD}
export CF_API_URL=${CF_API_URL}
EOL
