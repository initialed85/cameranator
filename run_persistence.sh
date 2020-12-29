#!/usr/bin/env bash

set -e -x

pushd "$(pwd)"

function teardown() {
  docker-compose down --remove-orphans --volumes || true
  popd >/dev/null 2>&1
}
trap teardown EXIT

export DOCKER_BUILDKIT=1

docker-compose up -d postgres hasura

export HASURA_GRAPHQL_ENDPOINT="http://host.docker.internal:8079/"

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

docker-compose logs -f -t
