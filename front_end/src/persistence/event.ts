import moment from "moment/moment";
import { getVideo, Video } from "./video";
import { getImage, Image } from "./image";
import { DocumentNode, gql } from "@apollo/client";
import { getClient } from "./utils";
import { Camera, getCamera } from "./camera";

export interface Event {
    uuid: string;
    start_timestamp: moment.Moment;
    end_timestamp: moment.Moment;
    is_segment: boolean;
    high_quality_video: Video;
    high_quality_image: Image;
    low_quality_video: Video;
    low_quality_image: Image;
    source_camera: Camera;
}

export interface GetEventsHandler {
    (events: Event[]): void;
}

function getQuery(
    isSegment: boolean,
    date: moment.Moment,
    cameraName: string
): DocumentNode {
    const startTimestamp = moment(
        `${date.local().format("YYYY-MM-DD")}T00:00:00+0800`
    );
    const endTimestamp = moment(
        `${date.local().format("YYYY-MM-DD")}T23:59:00+0800`
    );

    return gql(`
{
  event(
    order_by: {start_timestamp: desc}, 
    where: {
      is_segment: {_eq: ${isSegment}}, 
      start_timestamp: {_gte: "${startTimestamp.toISOString()}"}, 
      end_timestamp: {_lte: "${endTimestamp.toISOString()}"}
      source_camera: {name: {_eq: "${cameraName}"}}
    }
  ) {
    uuid
    start_timestamp
    end_timestamp
    is_segment
    high_quality_video {
      uuid
      file_path
    }
    high_quality_image {
      uuid
      file_path
    }
    low_quality_video {
      uuid
      file_path
    }
    low_quality_image {
      uuid
      file_path
    }
    source_camera {
      uuid
      name
    }
  }
}
`);
}

export function getEvent(item: any): Event {
    return {
        uuid: item["uuid"],
        start_timestamp: moment.utc(item["start_timestamp"]),
        end_timestamp: moment.utc(item["end_timestamp"]),
        is_segment: item["is_segment"],
        high_quality_video: getVideo(item["high_quality_video"]),
        high_quality_image: getImage(item["high_quality_image"]),
        low_quality_video: getVideo(item["low_quality_video"]),
        low_quality_image: getImage(item["low_quality_image"]),
        source_camera: getCamera(item["source_camera"]),
    };
}

export function getEvents(
    isSegment: boolean,
    date: moment.Moment,
    cameraName: string,
    handler: GetEventsHandler
) {
    const client = getClient();

    client
        .query({
            query: getQuery(isSegment, date, cameraName),
        })
        .catch((e) => {
            console.error(e);
        })
        .then((r) => {
            let events: Event[] = [];

            const data = (r as any).data["event"].slice();

            data.forEach((item: any) => {
                events.push(getEvent(item));
            });

            handler(events);
        });
}
