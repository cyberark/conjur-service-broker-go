#!/bin/sh

set -e

conjur init -u http://conjur -a dev --insecure

conjur login -i "$CONJUR_AUTHN_LOGIN" -p "$CONJUR_AUTHN_API_KEY"

conjur whoami

sleep infinity

set +e
