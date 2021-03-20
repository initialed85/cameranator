FROM linuxserver/ffmpeg:4.3-cli-ls14 AS base

RUN apt-get update && apt-get upgrade -y libfontconfig1 && apt-get install -y --reinstall \
    libfontconfig1 libfontconfig1-dev fontconfig-config

FROM base AS build

RUN apt-get update && apt-get install -y \
    git autoconf automake pkgconf libtool build-essential libzip-dev gettext autopoint \
    libmicrohttpd-dev libjpeg8-dev libavutil-dev libavformat-dev libavcodec-dev libswscale-dev libavdevice-dev

WORKDIR /srv/
RUN git clone https://github.com/Motion-Project/motion.git && cd motion && git checkout 'release-4.3.2' && cd ..

WORKDIR /srv/motion
RUN autoreconf -fiv
RUN ./configure --prefix=/
RUN make
RUN make install

FROM base AS run

ENV TZ Australia/Perth
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && apt-get install -y \
    tzdata netcat net-tools inetutils-ping \
    x264 netcat libmicrohttpd-dev libjpeg8-dev libavutil-dev libavformat-dev libavcodec-dev libswscale-dev libavdevice-dev

RUN dpkg-reconfigure -f noninteractive tzdata

# motion web ui
EXPOSE 8080

# motion streams
EXPOSE 8081

COPY --from=build /bin/motion /bin/motion
COPY --from=build /etc/motion/motion-dist.conf /etc/motion/examples/motion.conf
COPY --from=build /etc/motion/camera1-dist.conf /etc/motion/examples/camera1.conf
COPY --from=build /etc/motion/camera2-dist.conf /etc/motion/examples/camera2.conf
COPY --from=build /etc/motion/camera3-dist.conf /etc/motion/examples/camera3.conf

RUN mkdir -p /etc/motion/conf.d

WORKDIR /srv/

COPY docker/motion/handle_event.sh /srv/handle_event.sh

ENV UDP_HOST localhost
ENV UDP_PORT 6291

ENTRYPOINT ["motion", "-n"]
