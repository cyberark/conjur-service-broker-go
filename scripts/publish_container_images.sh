#!/bin/bash

set -e

# Navigate to the bin directory (where this script lives) to ensure we can run this script
# from anywhere.
cd "$(dirname "${0}")"

. build_utils.sh

function print_help() {
	echo "Internal Release Usage: ${0} --internal"
	echo "External Release Usage: ${0}"
	echo "Promote Usage: ${0} --promote --source <VERSION> --target <VERSION>"
	echo " --internal: publish images to registry.tld"
	echo " --source <VERSION>: specify version number of local image"
	echo " --target <VERSION>: specify version number of remote image"
}

# Fail if no arguments are given.
if [[ $# -lt 1 ]]; then
	print_help
	exit 1
fi

PUBLISH_INTERNAL=false
PROMOTE=false
PULL_SOURCE_IMAGES=false

while [[ $# -gt 0 ]]; do
	case "${1}" in
	--internal)
		PUBLISH_INTERNAL=true
		;;
	--pull)
		PULL_SOURCE_IMAGES=true
		;;
	--promote)
		PROMOTE=true
		;;
	--source)
		SOURCE_ARG="${2}"
		shift
		;;
	--target)
		TARGET_ARG="${2}"
		shift
		;;
	--help)
		print_help
		exit 1
		;;
	*)
		echo "Unknown option: ${1}"
		print_help
		exit 1
		;;
	esac
	shift
done

readonly INTERNAL_REGISTRY="registry.tld"
# Version derived from CHANGLEOG and automated release library
VERSION_WITH_COMMIT="$(project_version_with_commit)"
readonly VERSION_WITH_COMMIT
readonly IMAGE_NAME="conjur-service-broker"


if [[ ${PUBLISH_INTERNAL} = true ]]; then
	echo "Publishing built images internally to registry.tld."
	SOURCE_TAG=${VERSION_WITH_COMMIT}
	REMOTE_TAG=${VERSION_WITH_COMMIT}

	echo "SOURCE_TAG=${SOURCE_TAG}, REMOTE_TAG=${REMOTE_TAG}"
	tag_and_push "${IMAGE_NAME}:${SOURCE_TAG}" "${INTERNAL_REGISTRY}/${IMAGE_NAME}:${REMOTE_TAG}"
fi

if [[ ${PROMOTE} = true ]]; then
	if [[ -z ${SOURCE_ARG:-} || -z ${TARGET_ARG:-} ]]; then
		echo "When promoting, --source and --target flags are required."
		print_help
		exit 1
	fi

	# Update vars to utilize build_utils.sh
	SOURCE_TAG=${SOURCE_ARG}
	REMOTE_TAG=${TARGET_ARG}

	echo "Promoting image to ${REMOTE_TAG}"
	readonly TAGS=(
		"${REMOTE_TAG}"
		"latest"
	)

	if [[ ${PULL_SOURCE_IMAGES} = true ]]; then
		echo "Pulling source images from local registry"
		docker pull "${INTERNAL_REGISTRY}/${IMAGE_NAME}:${SOURCE_TAG}"
	fi

	for tag in "${TAGS[@]}" $(gen_versions "${REMOTE_TAG}"); do
		tag_and_push "${INTERNAL_REGISTRY}/${IMAGE_NAME}:${SOURCE_TAG}" "${INTERNAL_REGISTRY}/${IMAGE_NAME}:${tag}"
	done
fi
