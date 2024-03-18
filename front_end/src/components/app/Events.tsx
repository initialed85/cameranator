import { Camera } from "../../hasura/camera"
import moment from "moment"
import { Type } from "../../hasura/type"
import { useSubscription } from "@apollo/client"
import { Detection, Event, getEventsQuery } from "../../hasura/event"
import { Table } from "react-bootstrap"
import { useState } from "react"
import { adjustPath, Preview } from "./Preview"
import { CloudDownload } from "react-bootstrap-icons"

const MIN_SECONDS_SEEN = 2.0
const MIN_SCORE = 0.55
const TOP_N_SCORES = 5

export interface EventsProps {
    responsive: boolean
    camera: Camera | null
    date: moment.Moment | null
    type: Type | null
    objectFilter: string
}

export interface DetectionSummary {
    className: string
    weightedScore: number
    score: number
}

export function Events(props: EventsProps) {
    const eventsQuery = useSubscription(
        getEventsQuery(props.camera!, props.date!, props.type!),
    )

    const [focusEventUUID, setFocusEventUUID] = useState(null)

    const rows: JSX.Element[] = []

    const detectionsByEventID: Map<string, Detection[]> = new Map()

    const rawEvents = eventsQuery?.data?.event_with_detection || []

    rawEvents.forEach((event: Event) => {
        if (!event.class_name) {
            return
        }

        const detections: Detection[] = detectionsByEventID.get(event.id) || []

        detections.push({
            class_id: event.class_id,
            class_name: event.class_name,
            score: event.score,
            count: event.count,
            weighted_score: event.weighted_score,
        })

        detectionsByEventID.set(event.id, detections)
    })

    const eventByID: Map<string, Event> = new Map()
    rawEvents.forEach((event: Event) => {
        event.detections = detectionsByEventID.get(event.id) || []

        eventByID.set(event.id, event)
    })

    let events: Event[] = []
    eventByID.forEach((event: Event) => {
        events.push(event)
    })

    events = events.sort((a, b) => {
        if (a.start_timestamp > b.start_timestamp) {
            return -1
        } else if (a.start_timestamp < b.start_timestamp) {
            return 1
        }

        return 0
    })

    events.forEach((event: Event) => {
        const startTimestamp = moment.utc(event.start_timestamp)
        const endTimestamp = moment.utc(event.end_timestamp)

        const duration = moment.duration(
            endTimestamp.diff(event.start_timestamp),
        )

        const objectFilter = (props.objectFilter.trim() || "").toLowerCase()

        const detectionByClassName: Map<string, Detection> = new Map()

        const detections = event.detections || []
        detections.forEach((detection) => {
            if (objectFilter) {
                if (!detection?.class_name?.includes(objectFilter)) {
                    return
                }
            }

            detectionByClassName.set(detection.class_name, detection)
        })

        let detectionSummaries: DetectionSummary[] = []
        detectionByClassName.forEach((detection, className) => {
            // 20 fps / 4 stride frames = seconds seen
            if (detection.count / (20 / 4) < MIN_SECONDS_SEEN) {
                return
            }

            if (detection.score < MIN_SCORE) {
                return
            }

            if (detectionSummaries.length >= TOP_N_SCORES) {
                return
            }

            detectionSummaries.push({
                className,
                weightedScore: detection.weighted_score,
                score: detection.score,
            })
        })

        detectionSummaries = detectionSummaries.sort((a, b) => {
            if (a.weightedScore < b.weightedScore) {
                return 1
            } else if (a.weightedScore > b.weightedScore) {
                return -1
            }

            return 0
        })

        const objectElements: JSX.Element[] = []

        detectionSummaries.forEach((detectionSummary) => {
            objectElements.push(
                <div
                    style={{
                        display: "flex",
                        flexDirection: "row",
                        justifyContent: "space-between",
                    }}
                >
                    <span>{detectionSummary.className}</span>
                    <span>{detectionSummary.score.toFixed(2)}</span>
                </div>,
            )
        })

        if (objectFilter.length && !objectElements.length) {
            return
        }

        rows.push(
            <tr key={event.id}>
                <td style={{ verticalAlign: "middle" }}>
                    <div
                        style={{
                            display: "flex",
                            flexDirection: "column",
                        }}
                    >
                        {props.responsive ? (
                            <>
                                <span>
                                    {startTimestamp.local().format("HH:mm:ss")}
                                </span>
                                <span>
                                    {endTimestamp.local().format("HH:mm:ss")}
                                </span>
                                <span style={{ color: "gray" }}>
                                    {duration.minutes()}m{duration.seconds()}s
                                </span>
                                <span style={{ color: "gray" }}>
                                    {event.original_video.size.toFixed(0)} MB
                                </span>
                            </>
                        ) : (
                            <>
                                <span>
                                    {startTimestamp.local().format("HH:mm:ss")}{" "}
                                    to {endTimestamp.local().format("HH:mm:ss")}{" "}
                                    <span style={{ color: "gray" }}>
                                        ({duration.minutes()}m
                                        {duration.seconds()}s @{" "}
                                        {event.original_video.size.toFixed(0)}{" "}
                                        MB)
                                    </span>
                                </span>
                            </>
                        )}
                    </div>
                </td>
                <td style={{ verticalAlign: "middle", width: "125px" }}>
                    <div
                        style={{
                            display: "flex",
                            flexDirection: "column",
                            paddingLeft: "5px",
                            paddingRight: "5px",
                        }}
                    >
                        {objectElements}
                    </div>
                </td>
                <td
                    style={{
                        verticalAlign: "middle",
                        width: !props.responsive ? 330 : 130,
                    }}
                    onMouseOver={(e) => {
                        if (props.responsive) {
                            return
                        }

                        setFocusEventUUID(event.id as any)
                    }}
                    onMouseOut={(e) => {
                        if (props.responsive) {
                            return
                        }

                        setFocusEventUUID(null)
                    }}
                    onClick={(e) => {
                        if (props.responsive) {
                            return
                        }

                        setFocusEventUUID(null)
                    }}
                >
                    <a
                        target={`_thumbnail_image_${event.id}`}
                        rel={"noreferrer"}
                        href={adjustPath(event.thumbnail_image?.file_path)}
                    >
                        <Preview
                            event={event}
                            focusEventUUID={focusEventUUID}
                        />
                    </a>
                </td>
                <td style={{ verticalAlign: "middle" }}>
                    <a
                        target={`_original_video_${event.id}`}
                        rel={"noreferrer"}
                        href={adjustPath(event.original_video?.file_path)}
                    >
                        <CloudDownload
                            style={{
                                width: "20px",
                                height: "20px",
                                color: "gray",
                            }}
                        />
                    </a>
                </td>
            </tr>,
        )
    })

    return (
        <Table striped bordered hover size="sm">
            <thead>
                <tr>
                    <th>Details</th>
                    <th>Objects</th>
                    <th colSpan={2}>Media</th>
                </tr>
            </thead>
            <tbody>{rows}</tbody>
        </Table>
    )
}
