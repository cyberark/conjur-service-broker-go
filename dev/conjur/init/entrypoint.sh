#!/bin/sh

conjur authn login

conjur policy load service-broker.yaml


sleep infinity