import moment from "moment/moment";
import { Camera, getCamera } from "../types/camera";
import { DocumentNode, gql } from "@apollo/client";
import { getVideo } from "../types/video";
import { getImage } from "../types/image";
import Collection from "../collection";

export interface Video {
    uuid: string;
    file_path: string;
}

export interface Image {
    uuid: string;
    file_path: string;
}

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

function getQuery(args: any): DocumentNode {
    const startTimestamp = moment(
        `${args.date.local().format("YYYY-MM-DD")}T00:00:00+0800`
    );
    const endTimestamp = moment(
        `${args.date.local().format("YYYY-MM-DD")}T23:59:00+0800`
    );

    return gql(`
{
  event(
    order_by: {start_timestamp: desc}, 
    where: {
      is_segment: {_eq: ${args.isSegment}}, 
      start_timestamp: {_gte: "${startTimestamp.toISOString()}"}, 
      end_timestamp: {_lte: "${endTimestamp.toISOString()}"}
      source_camera: {name: {_eq: "${args.cameraName}"}}
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

export class EventCollection extends Collection {
    constructor() {
        super(getQuery, "event");
    }

    get(args: any): Promise<any> {
        return new Promise((resolve, reject) => {
            this.handleResultPromise(this.getResultPromise(args))
                .catch((e) => {
                    reject(e);
                })
                .then((data) => {
                    let events: Event[] = [];

                    data.forEach((item: any) => {
                        events.push(getEvent(item));
                    });

                    resolve(events);
                });
        });
    }
}
