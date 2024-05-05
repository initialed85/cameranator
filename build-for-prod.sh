#!/usr/bin/env bash

set -e -m

export GOOS=linux
export GOARCH=amd64

# if [[ "${1}" == "" || "${1}" == "quotanizer" ]]; then
#     docker build --platform=linux/amd64 -t initialed85/cameranator-quotanizer:latest -f quotanizer/Dockerfile quotanizer # &
# fi

# if [[ "${1}" == "" || "${1}" == "segment-processor" ]]; then
#     go build -v -o segment_processor ./cmd/segment_processor/main.go
#     docker build --platform=linux/amd64 -t initialed85/cameranator-segment-processor:latest -f docker/segment-processor/Dockerfile . # &
#     rm -f segment_processor
# fi

# if [[ "${1}" == "" || "${1}" == "segment-generator" ]]; then
#     go build -v -o segment_generator ./cmd/segment_generator/main.go
#     docker build --platform=linux/amd64 -t initialed85/cameranator-segment-generator:latest -f docker/segment-generator/Dockerfile . # &
#     rm -f segment_generator
# fi

# if [[ "${1}" == "" || "${1}" == "event-pruner" ]]; then
#     go build -v -o event_pruner ./cmd/event_pruner/main.go
#     docker build --platform=linux/amd64 -t initialed85/cameranator-event-pruner:latest -f docker/event-pruner/Dockerfile . # &
#     rm -f event_pruner
# fi

# if [[ "${1}" == "" || "${1}" == "front-end" ]]; then
#     docker build --platform=linux/amd64 -t initialed85/cameranator-front-end:latest -f docker/front-end/Dockerfile . # &
# fi

# if [[ "${1}" == "" || "${1}" == "object-task-scheduler" ]]; then
#     go build -v -o object_task_scheduler ./cmd/object_task_scheduler/main.go
#     docker build --platform=linux/amd64 -t kube-registry:5000/cameranator-object-task-scheduler:latest -f docker/object-task-scheduler/Dockerfile . # &
#     rm -f object_task_scheduler
# fi

if [[ "${1}" == "" || "${1}" == "object-task-worker" ]]; then
    docker build --platform=linux/amd64 -t kube-registry:5000/cameranator-object-task-worker-nvidia-generic:latest -f docker/object-task-worker/Dockerfile.nvidia-generic . # &
    # docker build --platform=linux/amd64 -t kube-registry:5000/cameranator-object-task-worker-nvidia-sm30:latest -f docker/object-task-worker/Dockerfile.nvidia-sm30 .       # &
    # docker build --platform=linux/amd64 -t kube-registry:5000/cameranator-object-task-worker-amd-generic:latest -f docker/object-task-worker/Dockerfile.amd-generic .       # &
fi

wait

# if [[ "${1}" == "" || "${1}" == "quotanizer" ]]; then
#     docker push initialed85/cameranator-quotanizer:latest # &
#     sleep 1
# fi

# if [[ "${1}" == "" || "${1}" == "segment-processor" ]]; then
#     docker push initialed85/cameranator-segment-processor:latest # &
#     sleep 1
# fi

# if [[ "${1}" == "" || "${1}" == "segment-generator" ]]; then
#     docker push initialed85/cameranator-segment-generator:latest # &
#     sleep 1
# fi

# if [[ "${1}" == "" || "${1}" == "event-pruner" ]]; then
#     docker push initialed85/cameranator-event-pruner:latest # &
#     sleep 1
# fi

# if [[ "${1}" == "" || "${1}" == "front-end" ]]; then
#     docker push initialed85/cameranator-front-end:latest # &
#     sleep 1
# fi

# if [[ "${1}" == "" || "${1}" == "object-task-scheduler" ]]; then
#     docker push kube-registry:5000/cameranator-object-task-scheduler:latest # &
#     sleep 1
# fi

if [[ "${1}" == "" || "${1}" == "object-task-worker" ]]; then
    docker push kube-registry:5000/cameranator-object-task-worker-nvidia-generic:latest # &
    sleep 1
    # docker push kube-registry:5000/cameranator-object-task-worker-nvidia-sm30:latest    # &
    # sleep 1
    # docker push kube-registry:5000/cameranator-object-task-worker-amd-generic:latest    # &
    # sleep 1
fi

# wait

# if [[ "${1}" == "" || "${1}" == "quotanizer" ]]; then
#     kubectl --context home -n cameranator rollout restart statefulset/quotanizer
# fi

# if [[ "${1}" == "" || "${1}" == "segment-processor" ]]; then
#     kubectl --context home -n cameranator rollout restart statefulset/segment
# fi

# if [[ "${1}" == "" || "${1}" == "segment-generator" ]]; then
#     kubectl --context home -n cameranator rollout restart statefulset/segment
# fi

# if [[ "${1}" == "" || "${1}" == "event-pruner" ]]; then
#     kubectl --context home -n cameranator rollout restart statefulset/pruner
# fi

# if [[ "${1}" == "" || "${1}" == "front-end" ]]; then
#     kubectl --context home -n cameranator rollout restart deployment/nginx
# fi

# if [[ "${1}" == "" || "${1}" == "object-task-scheduler" ]]; then
#     kubectl --context home -n cameranator rollout restart statefulset/object-task-scheduler
# fi

# if [[ "${1}" == "" || "${1}" == "object-task-worker" ]]; then
#     kubectl --context home -n cameranator rollout restart statefulset/object-task-worker-nvidia-generic
# fi
