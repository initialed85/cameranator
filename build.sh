#!/usr/bin/env bash

set -e -x

pushd "$(pwd)"
function teardown() {
  popd
}
trap teardown exit

export DOCKER_BUILDKIT=1

if [[ ! -d quotanizer ]]; then
  git clone https://github.com/initialed85/quotanizer.git
fi

cd quotanizer
git reset --hard
git pull --all
cd ..

if ! docker-compose build --parallel; then
  echo "error: build failed"
  exit 1
fi
