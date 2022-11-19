import { gql } from "@apollo/client";
import moment from "moment";
import { Camera } from "./camera";
import { Type } from "./type";

export const getEventsQuery = (
  camera: Camera,
  date: moment.Moment,
  type: Type
) => {
  if (type.is_stream) {
    throw Error(
      `getEventsQuery invoked w/ type.is_stream=true! (camera=${JSON.stringify(
        camera
      )}, date=${date?.toISOString()}, type=${JSON.stringify(type)})`
    );
  }

  const startTimestamp = moment(
    `${date.local().format("YYYY-MM-DD")}T00:00:00+0800`
  );
  const endTimestamp = moment(
    `${date.local().format("YYYY-MM-DD")}T23:59:00+0800`
  );

  return gql`
        query {
            event(
                order_by: {start_timestamp: desc},
                where: {
                    is_segment: {_eq: ${type.is_segment}},
                    start_timestamp: {_gte: "${startTimestamp.toISOString()}"},
                    end_timestamp: {_lte: "${endTimestamp.toISOString()}"}
                    source_camera: {uuid: {_eq: "${camera.uuid}"}}
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
    `;
};

export interface Video {
  uuid: string;
  file_path: string;
}

export interface Image {
  uuid: string;
  file_path: string;
}

export interface Event {
  __typename: string;
  uuid: string;
  start_timestamp: string;
  end_timestamp: string;
  is_segment: boolean;
  high_quality_video: Video;
  high_quality_image: Image;
  low_quality_video: Video;
  low_quality_image: Image;
  source_camera: Camera;
}
