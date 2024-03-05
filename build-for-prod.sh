#!/usr/bin/env bash

set -e -m

# docker build --platform=linux/amd64 \
#   -t initialed85/cameranator-quotanizer:latest \
#   -f quotanizer/Dockerfile \
#   quotanizer &
# sleep 0.1

# docker build --platform=linux/amd64 \
#   -t initialed85/cameranator-motion-processor:latest \
#   -f docker/motion-processor/Dockerfile \
#   . &
# sleep 0.1

# docker build --platform=linux/amd64 \
#   -t initialed85/cameranator-motion:latest \
#   -f docker/motion/Dockerfile \
#   . &
# sleep 0.1

# docker build --platform=linux/amd64 \
#   -t initialed85/cameranator-segment-processor:latest \
#   -f docker/segment-processor/Dockerfile \
#   . &
# sleep 0.1

# docker build --platform=linux/amd64 \
#   -t initialed85/cameranator-segment-generator:latest \
#   -f docker/segment-generator/Dockerfile \
#   . &
# sleep 0.1

# docker build --platform=linux/amd64 \
#   -t initialed85/cameranator-event-pruner:latest \
#   -f docker/event-pruner/Dockerfile \
#   . &
# sleep 0.1

docker build --platform=linux/amd64 \
  -t initialed85/cameranator-front-end:latest \
  -f docker/front-end/Dockerfile \
  . &
sleep 0.1

# docker build --platform=linux/amd64 \
#   -t kube-registry:5000/cameranator-object-task-scheduler:latest \
#   -f docker/object-task-scheduler/Dockerfile \
#   . &
# sleep 0.1

# docker build --platform=linux/amd64 \
#   -t kube-registry:5000/cameranator-object-task-worker-nvidia-sm30:latest \
#   -f docker/object-task-worker/Dockerfile.nvidia-sm30 \
#   . &
# sleep 0.1

# docker build --platform=linux/amd64 \
#   -t kube-registry:5000/cameranator-object-task-worker-nvidia-generic:latest \
#   -f docker/object-task-worker/Dockerfile.nvidia-generic \
#   . &
# sleep 0.1

# docker build --platform=linux/amd64 \
#   -t kube-registry:5000/cameranator-object-task-worker-amd-generic:latest \
#   -f docker/object-task-worker/Dockerfile.amd-generic \
#   . &
# sleep 0.1

wait

# docker push initialed85/cameranator-quotanizer:latest &
# sleep 0.1

# docker push initialed85/cameranator-motion-processor:latest &
# sleep 0.1

# docker push initialed85/cameranator-motion:latest &
# sleep 0.1

# docker push initialed85/cameranator-segment-processor:latest &
# sleep 0.1

# docker push initialed85/cameranator-segment-generator:latest &
# sleep 0.1

# docker push initialed85/cameranator-event-pruner:latest &
# sleep 0.1

docker push initialed85/cameranator-front-end:latest &
sleep 0.1

# docker push kube-registry:5000/cameranator-object-task-scheduler:latest &
# sleep 0.1

# docker push kube-registry:5000/cameranator-object-task-worker-nvidia-sm30:latest &
# sleep 0.1

# docker push kube-registry:5000/cameranator-object-task-worker-nvidia-generic:latest &
# sleep 0.1

# docker push kube-registry:5000/cameranator-object-task-worker-amd-generic:latest &
# sleep 0.1

wait
