import { Camera } from "../../hasura/camera"
import moment from "moment"
import { Type } from "../../hasura/type"
import { useSubscription } from "@apollo/client"
import { Detection, Event, getEventsQuery } from "../../hasura/event"
import { Table } from "react-bootstrap"
import { useEffect, useState } from "react"
import { adjustPath, Preview } from "./Preview"
import { CloudDownload } from "react-bootstrap-icons"

const MIN_SECONDS_SEEN = 2.0
// const MIN_SCORE = 0.55
const TOP_N_SCORES = 5

export interface EventsProps {
    responsive: boolean
    camera: Camera | null
    date: moment.Moment | null
    type: Type | null
    objectFilter: string
    setLoading: (x: boolean) => void
}

export interface DetectionSummary {
    className: string
    weightedScore: number
    score: number
    count: number
}

export function Events(props: EventsProps) {
    const eventsQuery = useSubscription(
        getEventsQuery(props.camera!, props.date!, props.type!),
    )

    const [focusEventUUID, setFocusEventUUID] = useState(null)

    const rows: JSX.Element[] = []

    const detectionsByEventID: Map<string, Detection[]> = new Map()

    const rawEvents = eventsQuery?.data?.event || []

    rawEvents.forEach((event: Event) => {
        const detections: Detection[] = detectionsByEventID.get(event.id) || []

        event.aggregated_detections.forEach((detection) => {
            detections.push(detection)
        })

        detectionsByEventID.set(event.id, detections)
    })

    const eventByID: Map<string, Event> = new Map()
    rawEvents.forEach((event: Event) => {
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

        const detections = event.aggregated_detections || []
        detections.forEach((detection) => {
            detectionByClassName.set(detection.class_name, detection)
        })

        const filteredDetections = detections.filter((detection) => {
            if (objectFilter) {
                if (!detection?.class_name?.includes(objectFilter)) {
                    return false
                }
            }

            return true
        })

        if (detections.length && !filteredDetections.length) {
            return
        }

        let detectionSummaries: DetectionSummary[] = []
        detectionByClassName.forEach((detection, className) => {
            // 20 fps / 4 stride frames = seconds seen
            if (detection.count / (20 / 4) < MIN_SECONDS_SEEN) {
                return
            }

            // if (detection.score < MIN_SCORE) {
            //     return
            // }

            if (detectionSummaries.length >= TOP_N_SCORES) {
                return
            }

            detectionSummaries.push({
                className,
                weightedScore: detection.weighted_score,
                score: detection.score,
                count: detection.count,
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
            const matchesObjectFilter =
                objectFilter.length &&
                detectionSummary.className.includes(objectFilter)

            const color = matchesObjectFilter ? "red" : "black"

            objectElements.push(
                <div
                    style={{
                        display: "flex",
                        flexDirection: "row",
                        justifyContent: "space-between",
                    }}
                >
                    <span
                        style={{
                            color,
                        }}
                    >
                        {detectionSummary.className}
                    </span>
                    <span
                        style={{
                            color,
                        }}
                    >
                        {detectionSummary.count} @{" "}
                        {detectionSummary.score.toFixed(2)}
                    </span>
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
                <td style={{ verticalAlign: "middle", width: "200px" }}>
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
                    <div
                        style={{
                            display: "flex",
                            flexDirection: "row",
                            justifyContent: "center",
                        }}
                    >
                        <></>
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
                        <></>
                    </div>
                </td>
                <td style={{ verticalAlign: "middle", width: "50px" }}>
                    <a
                        target={`_original_video_${event.id}`}
                        rel={"noreferrer"}
                        href={adjustPath(event.original_video?.file_path)}
                    >
                        <CloudDownload
                            style={{
                                width: "25px",
                                height: "25px",
                                color: "gray",
                            }}
                        />
                    </a>
                </td>
            </tr>,
        )
    })

    useEffect(() => {
        props.setLoading(eventsQuery.loading)
    }, [eventsQuery.loading, props])

    return (
        <Table striped bordered hover size="sm">
            <thead>
                <tr>
                    <th>Details</th>
                    <th>Detections</th>
                    <th colSpan={2}>Media</th>
                </tr>
            </thead>
            <tbody>{rows}</tbody>
        </Table>
    )
}
