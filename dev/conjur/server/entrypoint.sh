#!/bin/sh

set -e

psql "$DATABASE_URL" -c "DROP SCHEMA IF EXISTS public CASCADE; CREATE SCHEMA public;"
#export CONJUR_DATA_KEY=$(conjurctl data-key generate)
#echo "export CONJUR_DATA_KEY=${CONJUR_DATA_KEY}" > /etc/profile.d/02-data-key.sh
#chmod +x /etc/profile.d/02-data-key.sh

conjurctl server &
conjurctl wait

conjurctl account create dev || true >/dev/null 2>&1

conjurctl policy load dev /init/service-broker.yaml || true >/dev/null 2>&1

touch up

wait

set +e
