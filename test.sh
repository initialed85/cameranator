#!/usr/bin/env bash

set -e -x

pushd "$(pwd)"

function teardown() {
  docker rm -f ffmpeg || true
  docker rm -f rtsp-simple-server || true
  docker-compose down --remove-orphans --volumes || true
  popd >/dev/null 2>&1 || true
}
trap teardown exit

CCTV_EVENTS_QUOTA=1
CCTV_EVENTS_PATH="$(pwd)/temp_data/events"

CCTV_SEGMENTS_QUOTA=1
CCTV_SEGMENTS_PATH="$(pwd)/temp_data/segments"
CCTV_SEGMENT_DURATION=30

CCTV_MOTION_CONFIGS="$(pwd)/motion-configs-testing"

DISABLE_NVIDIA=1

export CCTV_EVENTS_QUOTA
export CCTV_EVENTS_PATH

export CCTV_SEGMENTS_QUOTA
export CCTV_SEGMENTS_PATH
export CCTV_SEGMENT_DURATION

export CCTV_MOTION_CONFIGS

export DISABLE_NVIDIA

export DOCKER_BUILDKIT=1

export HASURA_GRAPHQL_ENDPOINT="http://localhost:8082/"

#
# hasura deps
#

docker-compose up -d postgres hasura

cd persistence/hasura

while ! hasura migrate apply; do
  sleep 1
done

hasura metadata apply

hasura seeds apply

popd

#
# test stream
#

docker run --rm -d --name rtsp-simple-server -e RTSP_PROTOCOLS=tcp -p 8554:8554 aler9/rtsp-simple-server

docker run --rm -d --name ffmpeg -v "$(pwd)/test_data/segments/":/srv/ jrottenberg/ffmpeg:4.3.1-ubuntu1804 \
  -re -stream_loop -1 -i /srv/Segment_2020-12-25T08:45:04_Driveway.mp4 -c copy -f rtsp rtsp://host.docker.internal:8554/Streaming/Channels/101

#
# run tests (serially, because of the shared database)
#

go test -v ./pkg/filesystem

go test -v ./pkg/front_end/legacy/page_render_utils
go test -v ./pkg/front_end/legacy/page_renderer

go test -v ./pkg/media/converter
go test -v ./pkg/media/metadata
go test -v ./pkg/media/segment_recorder
go test -v ./pkg/media/thumbnail_creator

go test -v ./pkg/motion/event_receiver

go test -v ./pkg/persistence/application
go test -v ./pkg/persistence/graphql
go test -v ./pkg/persistence/helpers
go test -v ./pkg/persistence/legacy
go test -v ./pkg/persistence/model
go test -v ./pkg/persistence/registry

go test -v ./pkg/process

go test -v ./pkg/segments/event_receiver
go test -v ./pkg/segments/segment_generator

go test -v ./pkg/services/motion_processor
go test -v ./pkg/services/page_renderers
go test -v ./pkg/services/segment_generators
go test -v ./pkg/services/segment_processor

go test -v ./pkg/test_utils
go test -v ./pkg/utils

echo ""
echo "All passed."
