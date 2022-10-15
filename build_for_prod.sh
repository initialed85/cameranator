#!/usr/bin/env bash

set -e -x

pushd "$(pwd)"

# docker build \
#   -t initialed85/cameranator-quotanizer \
#   -f quotanizer/Dockerfile \
#   quotanizer
#
# docker build \
#   -t initialed85/cameranator-motion-processor \
#   -f docker/motion-processor/Dockerfile \
#   .
#
# docker build \
#   -t initialed85/cameranator-motion \
#   -f docker/motion/Dockerfile \
#   .
#
# docker build \
#   -t initialed85/cameranator-segment-processor \
#   -f docker/segment-processor/Dockerfile \
#   .
#
# docker build \
#   -t initialed85/cameranator-segment-generator \
#   -f docker/segment-generator/Dockerfile \
#   .
#
# docker build \
#   -t initialed85/cameranator-event-pruner \
#   -f docker/event-pruner/Dockerfile \
#   .
#
docker build \
 -t initialed85/cameranator-front-end \
 -f docker/front-end/Dockerfile \
 .

#docker push initialed85/cameranator-quotanizer
#docker push initialed85/cameranator-motion-processor
#docker push initialed85/cameranator-motion
#docker push initialed85/cameranator-segment-processor
#docker push initialed85/cameranator-segment-generator
#docker push initialed85/cameranator-event-pruner
docker push initialed85/cameranator-front-end
