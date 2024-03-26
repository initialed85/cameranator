#!/usr/bin/env bash

set -e -m

export GOOS=linux
export GOARCH=amd64

# docker build --platform=linux/amd64 -t initialed85/cameranator-quotanizer:latest -f quotanizer/Dockerfile quotanizer # &

# go build -v -o segment_processor ./cmd/segment_processor/main.go
# docker build --platform=linux/amd64 -t initialed85/cameranator-segment-processor:latest -f docker/segment-processor/Dockerfile . # &
# rm -f segment_processor

# go build -v -o segment_generator ./cmd/segment_generator/main.go
# docker build --platform=linux/amd64 -t initialed85/cameranator-segment-generator:latest -f docker/segment-generator/Dockerfile . # &
# rm -f segment_generator

# go build -v -o event_pruner ./cmd/event_pruner/main.go
# docker build --platform=linux/amd64 -t initialed85/cameranator-event-pruner:latest -f docker/event-pruner/Dockerfile . # &
# rm -f event_pruner

docker build --platform=linux/amd64 -t initialed85/cameranator-front-end:latest -f docker/front-end/Dockerfile . # &

# go build -v -o object_task_scheduler ./cmd/object_task_scheduler/main.go
# docker build --platform=linux/amd64 -t kube-registry:5000/cameranator-object-task-scheduler:latest -f docker/object-task-scheduler/Dockerfile . # &
# rm -f object_task_scheduler

# docker build --platform=linux/amd64 -t kube-registry:5000/cameranator-object-task-worker-nvidia-generic:latest -f docker/object-task-worker/Dockerfile.nvidia-generic . # &

# docker build --platform=linux/amd64 -t kube-registry:5000/cameranator-object-task-worker-nvidia-sm30:latest -f docker/object-task-worker/Dockerfile.nvidia-sm30 .       # &

# docker build --platform=linux/amd64 -t kube-registry:5000/cameranator-object-task-worker-amd-generic:latest -f docker/object-task-worker/Dockerfile.amd-generic .       # &

wait

# docker push initialed85/cameranator-quotanizer:latest # &
# sleep 1

# docker push initialed85/cameranator-segment-processor:latest # &
# sleep 1

# docker push initialed85/cameranator-segment-generator:latest # &
# sleep 1

# docker push initialed85/cameranator-event-pruner:latest # &
# sleep 1

docker push initialed85/cameranator-front-end:latest # &
sleep 1

# docker push kube-registry:5000/cameranator-object-task-scheduler:latest # &
# sleep 1

# docker push kube-registry:5000/cameranator-object-task-worker-nvidia-generic:latest # &
# sleep 1

# docker push kube-registry:5000/cameranator-object-task-worker-nvidia-sm30:latest    # &
# sleep 1

# docker push kube-registry:5000/cameranator-object-task-worker-amd-generic:latest    # &
# sleep 1

wait
