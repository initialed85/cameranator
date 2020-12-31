#!/usr/bin/env bash

set -e -x

pushd "$(pwd)"

function teardown() {
  docker-compose down --remove-orphans || true
  popd >/dev/null 2>&1
}
trap teardown EXIT

CCTV_EVENTS_QUOTA=2000
CCTV_EVENTS_PATH="/media/storage/Cameras/events"

CCTV_SEGMENTS_QUOTA=4000
CCTV_SEGMENTS_PATH="/media/storage/Cameras/segments"
CCTV_SEGMENT_DURATION=300

CCTV_MOTION_CONFIGS="$(pwd)/motion-configs-run-in-prod"

CCTV_EXPOSE_PORT=81

export CCTV_EVENTS_QUOTA
export CCTV_EVENTS_PATH

export CCTV_SEGMENTS_QUOTA
export CCTV_SEGMENTS_PATH
export CCTV_SEGMENT_DURATION

export CCTV_MOTION_CONFIGS

export CCTV_EXPOSE_PORT

export DOCKER_BUILDKIT=1

export HASURA_GRAPHQL_ENDPOINT="http://localhost:8082/"

docker-compose up -d nginx postgres hasura motion

cd persistence/hasura

while ! hasura migrate apply; do
  sleep 1
done

while ! hasura metadata apply; do
  sleep 1
done

while ! hasura seeds apply; do
  sleep 1
done

popd

docker-compose up -d

docker-compose ps -a
