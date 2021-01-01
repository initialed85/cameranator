# cameranator

A ghetto CCTV system built on the following pieces of software:

- [motion](https://github.com/Motion-Project/motion)
- [FFmpeg](https://github.com/FFmpeg/FFmpeg)
- [ImageMagick](https://github.com/ImageMagick/ImageMagick)
- [Hasura](https://github.com/hasura) (for [GraphQL](https://graphql.org/) support)
- [Postgres](https://github.com/postgres/postgres)
- [quotanizer](https://github.com/initialed85/quotanizer)

Held together with the following languages / frameworks:
 
- [Go](https://github.com/golang)
- [glue](https://github.com/initialed85/glue)
- [TypeScript](https://github.com/microsoft/TypeScript) (TODO)
- [React](https://github.com/facebook/react) (TODO)

Deployed using the following pieces of software:  

- [Docker](https://github.com/docker/docker-ce) 
- [docker-compose](https://github.com/docker/compose).

## Concept

The system seeks to do the following:

- integrate with motion for motion detection
- integrate with FFmpeg for constant segment recording (TODO)
- create low resolution previews from captured videos and images
- write metadata to a database (using GraphQL)
- provide a web-based UI for the user (using TypeScript and React)

By having GraphQL at the middle of it all, the system _should_ be extensible in the future 
for things like WebSocket push events, other front-ends, mobile app, third-party integrations
etc.   

## Services

Each service below deploys as a Docker container:

- `postgres`
    - provide the database
- `hasura`
    - provide the schema, migrations and API
- `motion`
    - consume RTSP from cameras
    - detect motion events
    - generate `.mp4` and `.jpg` files for motion events
    - trigger shell scripts to send UDP event messages to `motion-processor`   
- `motion-processor`
    - consume UDP event messages from `motion`
    - convert the `.mp4` and `.jpg` files to low resolution previews (keeping the originals)
    - wrap up all the metadata and file paths as an event and push it to `hasura`
- `segment-recorder`
    - generate 5 minute `.mp4` segments using FFmpeg
    - extract `.jpg` files from the `.mp4` files
    - use file watchers to send UDP event messages to `segment-processor`
- `segment-processor`
    - consume UDP event messages from `segment-recorder`
    - convert the `.mp4` and `.jpg` files to low resolution previews (keeping the originals) 
    - wrap up all the metadata and file paths as an event and push it to `hasura`
- `page-renderer`
    - consume events from `hasura`
    - present them to the user
    - NOTE: legacy static page stuff
- `front-end` (TODO)
    - consume events from `hasura`
    - present them to the user

## Overall TODOs

- extend the persistence piece to be able to subscribe to push events
    - implied: work out how hasura does subscriptions (WebSockets?) 
- support timezones properly
- have a configuration system that extends out to motion and the other services
    - items to expose via config
        - disable nvidia
        - timezone for the system
        - nginx port
        - camera definitions
        - event quota
        - event path
        - segment quota
        - segment path
        - segment duration

## Technical debt

It wouldn't be a project (or at least a project that I wrote) without technical debt; here
are some things that need attention:

- the reflection / introspection part of the GraphQL query generation piece needs DRYing up
- various parts should have their own prefixed loggers (rather than just using `log.Printf`)
- the `Dockerfiles` for the services were split out from a larger monolithic `Dockerfile` and
  so there are probably some unrequired dependencies in them
- support for Nvidia is too baked in- it can only be disabled at a configuration level for the
  Go code; disabling it for `motion` requires configuration file changes (on the `motion` 
  side)
- camera config management isn't clean- one needs to define the necessary `motion` config 
  files in addition to the `camera` object instances in `hasura`
- there's a lot of repetition between the motion_processor and the segment_processor
- need to DRY up the Dockerfiles
- need to throw away the page render stuff (maybe not though? handy as an alternate path to see the data)

## Usage

### Building

```
./build.sh
```

This uses docker-compose to build all the Docker containers.

### Testing

```
./test.sh
```

This requires you to have run `./build.sh` first- it uses docker-compose to start the Docker
containers for the backing services (`postgres`, `hasura`) and some test / mock Docker
containers (`rtsp-simple-server`, `ffmpeg`) and finally runs the tests (which are a mixture
of unit and integration tests).

### Running locally

```
./run.sh
```

This requires you to have run `./build.sh` first- it uses docker-compose to start all the
Docker containers and it uses the configuration for my cameras at home from the 
`motion-config` folder (so it probably won't work for you without tweaking those configs).

On the assumption you've figured that piece out, navigate to [http://localhost/](http://localhost/)

### Other

There is also a `./run_persistence.sh` option- this is handy if you want to run the tests from
your IDE, and you're making changes to the persistence side of things (GraphQL query 
generation, etc).

Related to the above, you'll be able to access the Hasura admin UI 
[at this URL](http://localhost:8080/) when the above script is running.

If you're making changes that require an RTSP stream, you'll want to run the following (see
`./test.sh` for a reference):

```
# shell 1
docker run --rm -it --name rtsp-simple-server -e RTSP_PROTOCOLS=tcp -p 8554:8554 aler9/rtsp-simple-server

# shell 2
docker run --rm -it --name ffmpeg -v "$(pwd)/test_data/segments/":/srv/ jrottenberg/ffmpeg:4.3.1-ubuntu1804 \
  -re -stream_loop -1 -i /srv/Segment_2020-12-25T08:45:04_Driveway.mp4 -c copy -f rtsp rtsp://localhost:8554/Streaming/Channels/101
```

At this point, you'll have a looping RTSP stream of the folks who come and look after our 
cats when we're away [at this URL](rtsp://localhost:8554/Streaming/Channels/101).
