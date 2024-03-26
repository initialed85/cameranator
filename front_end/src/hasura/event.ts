import { gql } from "@apollo/client"
import moment from "moment"
import { Camera } from "./camera"
import { Type } from "./type"

export const getEventsQuery = (
    camera: Camera,
    date: moment.Moment,
    type: Type,
) => {
    if (type.is_stream) {
        throw Error(
            `getEventsQuery invoked w/ type.is_stream=true! (camera=${JSON.stringify(
                camera,
            )}, date=${date?.toISOString()}, type=${JSON.stringify(type)})`,
        )
    }

    const startTimestamp = moment(
        `${date.local().format("YYYY-MM-DD")}T00:00:00+0800`,
    )
    const endTimestamp = moment(
        `${date.local().format("YYYY-MM-DD")}T23:59:59+0800`,
    )

    return gql`
    subscription {
        event(
            order_by: {start_timestamp: desc},
            where: {
                start_timestamp: {_gte: "${startTimestamp.toISOString()}"},
                end_timestamp: {_lte: "${endTimestamp.toISOString()}"}
                source_camera: {id: {_eq: ${camera.id}}}
            }
        ) {
            id
            start_timestamp
            end_timestamp
            original_video {
                id
                file_path
                size
            }
            thumbnail_image {
                id
                file_path
            }
            processed_video {
                id
                file_path
            }
            source_camera {
                id
                name
            }
            aggregated_detections(order_by: {weighted_score: desc}) {
                class_id
                class_name
                score
                count
                weighted_score
            }
        }
    }`
}

export interface Video {
    id: string
    file_path: string
    size: number
}

export interface Image {
    id: string
    file_path: string
}

export interface Object {
    class_id: number
    class_name: string
    start_timestamp: string
    end_timestamp: string
    tracked_object_id: number
}

export interface Detection {
    class_id: number
    class_name: string
    score: number
    count: number
    weighted_score: number
}

export interface Event {
    __typename: string | null
    id: string
    start_timestamp: string
    end_timestamp: string
    original_video: Video
    thumbnail_image: Image
    processed_video: Video
    source_camera: Camera
    aggregated_detections: [Detection]
}
