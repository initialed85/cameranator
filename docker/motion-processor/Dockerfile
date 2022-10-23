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

RUN --mount=type=cache,target=/root/.cache/go-build go build -v -o motion_processor ./cmd/motion_processor/main.go

FROM base AS run

ENV TZ Australia/Perth
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y \
    tzdata imagemagick net-tools inetutils-ping \
    x264 libmicrohttpd-dev libjpeg8-dev libavutil-dev libavformat-dev libavcodec-dev libswscale-dev libavdevice-dev

RUN dpkg-reconfigure -f noninteractive tzdata

COPY --from=build /srv/motion_processor /srv/

WORKDIR /srv/

# motion event receiver
EXPOSE 6291/udp

ENTRYPOINT ["/srv/motion_processor"]

CMD []
