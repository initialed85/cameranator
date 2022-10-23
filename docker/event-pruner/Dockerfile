FROM linuxserver/ffmpeg AS base

RUN apt-get update && apt-get upgrade -y libfontconfig1 && apt-get install -y --reinstall \
    libfontconfig1 libfontconfig1-dev fontconfig-config

FROM golang:1.15 AS build

WORKDIR /srv/

COPY ./go.mod /srv/go.mod
COPY ./go.sum /srv/go.sum
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY ./cmd /srv/cmd
COPY ./pkg /srv/pkg
RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o event_pruner ./cmd/event_pruner/main.go

FROM base AS run

ENV TZ Australia/Perth
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y \
    tzdata

RUN dpkg-reconfigure -f noninteractive tzdata

COPY --from=build /srv/event_pruner /srv/

WORKDIR /srv/

# segment event receiver
EXPOSE 6291/udp

ENTRYPOINT ["/srv/event_pruner"]

CMD []
