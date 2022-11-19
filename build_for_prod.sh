#!/usr/bin/env bash

set -e -m

#docker build \
#  -t initialed85/cameranator-quotanizer:latest \
#  -f quotanizer/Dockerfile \
#  quotanizer &
#docker build \
#  -t initialed85/cameranator-motion-processor:latest \
#  -f docker/motion-processor/Dockerfile \
#  . &
#docker build \
#  -t initialed85/cameranator-motion:latest \
#  -f docker/motion/Dockerfile \
#  . &
#docker build \
#  -t initialed85/cameranator-segment-processor:latest \
#  -f docker/segment-processor/Dockerfile \
#  . &
#docker build \
#  -t initialed85/cameranator-segment-generator:latest \
#  -f docker/segment-generator/Dockerfile \
#  . &
#docker build \
#  -t initialed85/cameranator-event-pruner:latest \
#  -f docker/event-pruner/Dockerfile \
#  . &
docker build \
  -t initialed85/cameranator-front-end:latest \
  -f docker/front-end/Dockerfile \
  . &

wait

#docker push initialed85/cameranator-quotanizer:latest &
#docker push initialed85/cameranator-motion-processor:latest &
#docker push initialed85/cameranator-motion:latest &
#docker push initialed85/cameranator-segment-processor:latest &
#docker push initialed85/cameranator-segment-generator:latest &
#docker push initialed85/cameranator-event-pruner:latest &
docker push initialed85/cameranator-front-end:latest &

wait
