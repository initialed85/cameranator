FROM linuxserver/ffmpeg AS base

RUN apt-get update && apt-get upgrade -y libfontconfig1 && apt-get install -y --reinstall \
    libfontconfig1 libfontconfig1-dev fontconfig-config

FROM golang:1.15 AS build

WORKDIR /srv/

COPY ./go.mod /srv/go.mod
COPY ./go.sum /srv/go.sum
COPY ./cmd /srv/cmd
COPY ./pkg /srv/pkg

RUN go build -v -o segment_generator ./cmd/segment_generator/main.go

FROM base AS run

ENV TZ Australia/Perth
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y \
    tzdata imagemagick net-tools inetutils-ping \
    x264 libmicrohttpd-dev libjpeg8-dev libavutil-dev libavformat-dev libavcodec-dev libswscale-dev libavdevice-dev

RUN dpkg-reconfigure -f noninteractive tzdata

COPY --from=build /srv/segment_generator /srv/

WORKDIR /srv/

ENTRYPOINT ["/srv/segment_generator"]

CMD []
