import { gql } from "@apollo/client"

export const getDetectionsForEventQuery = (eventId: string) => {
    return gql`
        query {
            detections(
                where: { event_id: { _eq: ${eventId} } }
                order_by: { timestamp: asc }
            ) {
                timestamp
                class_id
                class_name
                score
                bounding_box
                centroid
            }
        }
    `
}

export interface Point {
    x: number
    y: number
}

export interface Detection {
    timestamp: string
    class_id: number
    class_name: string
    score: number
    count: number
    weighted_score: number
    bounding_box: any
    centroid: any
    timestampMilliseconds: number
    boundingBoxPoints: Point[]
    centroidPoint: Point
}
