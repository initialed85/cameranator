import { gql } from "@apollo/client"
import moment from "moment"
import { Camera } from "./camera"
import { Type } from "./type"
import { Detection } from "./detection"

export const getEventsSubscription = () => {
    return gql`
        subscription getEvents(
            $startTimestamp: timestamptz!
            $endTimestamp: timestamptz!
            $sourceCameraId: bigint!
        ) {
            event(
                limit: 30
                order_by: { start_timestamp: desc }
                where: {
                    start_timestamp: { _gte: $startTimestamp }
                    end_timestamp: { _lte: $endTimestamp }
                    source_camera: { id: { _eq: $sourceCameraId } }
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
                aggregated_detections(order_by: { weighted_score: desc }) {
                    class_id
                    class_name
                    score
                    count
                    weighted_score
                }
            }
        }
    `
}

export const getEventsQuery = () => {
    return gql`
        query getEvents(
            $startTimestamp: timestamptz!
            $endTimestamp: timestamptz!
            $sourceCameraId: bigint!
            $startTimestampLessThan: timestamptz!
        ) {
            event(
                limit: 30
                order_by: { start_timestamp: desc }
                where: {
                    start_timestamp: {
                        _gte: $startTimestamp
                        _lte: $startTimestampLessThan
                    }
                    end_timestamp: { _lte: $endTimestamp }
                    source_camera: { id: { _eq: $sourceCameraId } }
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
                aggregated_detections(order_by: { weighted_score: desc }) {
                    class_id
                    class_name
                    score
                    count
                    weighted_score
                }
            }
        }
    `
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
