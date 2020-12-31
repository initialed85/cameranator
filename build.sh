#!/usr/bin/env bash

set -e -x

pushd "$(pwd)"
function teardown() {
  popd
}
trap teardown exit

if [[ ! -d quotanizer ]]; then
  git clone https://github.com/initialed85/quotanizer.git
fi

cd quotanizer
git reset --hard
git pull --all
cd ..

CCTV_EVENTS_QUOTA=1
CCTV_EVENTS_PATH="$(pwd)/temp_data/events"

CCTV_SEGMENTS_QUOTA=1
CCTV_SEGMENTS_PATH="$(pwd)/temp_data/segments"
CCTV_SEGMENT_DURATION=30

CCTV_MOTION_CONFIGS="$(pwd)/motion-configs"

DISABLE_NVIDIA=1

export CCTV_EVENTS_QUOTA
export CCTV_EVENTS_PATH

export CCTV_SEGMENTS_QUOTA
export CCTV_SEGMENTS_PATH
export CCTV_SEGMENT_DURATION

export CCTV_MOTION_CONFIGS

export DISABLE_NVIDIA

export DOCKER_BUILDKIT=1

if ! docker-compose build --parallel; then
  echo "error: build failed"
  exit 1
fi
