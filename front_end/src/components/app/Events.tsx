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
const MIN_SCORE = 0.55
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
    outlier: boolean
}

export function Events(props: EventsProps) {
    const eventsQuery = useSubscription(
        getEventsQuery(props.camera!, props.date!, props.type!),
    )

    const [focusEventUUID, setFocusEventUUID] = useState(null)

    const rawObjectFilter = (props.objectFilter.trim() || "").replaceAll(
        " ",
        "",
    )
    const objectFilter =
        rawObjectFilter &&
        new RegExp(
            `^${(rawObjectFilter || "").toLowerCase().split(",").join("|")}$`,
        )

    const rows: JSX.Element[] = []

    let events: Event[] = eventsQuery?.data?.event || []

    events.forEach((event: Event) => {
        const startTimestamp = moment.utc(event.start_timestamp)
        const endTimestamp = moment.utc(event.end_timestamp)

        const duration = moment.duration(
            endTimestamp.diff(event.start_timestamp),
        )

        const detectionByClassName: Map<string, Detection> = new Map()
        const filteredDetectionByClassName: Map<string, Detection> = new Map()

        const detections = event.aggregated_detections || []
        detections.forEach((detection) => {
            detectionByClassName.set(detection.class_name, detection)

            if (objectFilter) {
                if (detection?.class_name?.match(objectFilter)) {
                    filteredDetectionByClassName.set(
                        detection.class_name,
                        detection,
                    )
                }
            }
        })

        if (
            objectFilter &&
            (!detectionByClassName.size ||
                (detectionByClassName.size &&
                    !filteredDetectionByClassName.size))
        ) {
            return
        }

        let detectionSummaries: DetectionSummary[] = []
        detectionByClassName.forEach((detection, className) => {
            // 20 fps / 4 stride frames = seconds seen
            const outlier =
                detection.count / (20 / 4) < MIN_SECONDS_SEEN ||
                detection.score < MIN_SCORE ||
                detectionSummaries.length >= TOP_N_SCORES

            detectionSummaries.push({
                className,
                weightedScore: detection.weighted_score,
                score: detection.score,
                count: detection.count,
                outlier: outlier,
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
                objectFilter && detectionSummary.className.match(objectFilter)

            let color = "black"
            let textDecoration = "none"
            if (detectionSummary.outlier) {
                color = "gray"
                textDecoration = "line-through"
            }

            if (matchesObjectFilter) {
                color = "red"
                if (detectionSummary.outlier) {
                    color = "maroon"
                }
            }

            let inner = (
                <>
                    <span
                        style={{
                            color,
                            textDecoration,
                        }}
                    >
                        {detectionSummary.className}
                    </span>
                    <span
                        style={{
                            color,
                            textDecoration,
                        }}
                    >
                        {detectionSummary.count} @{" "}
                        {detectionSummary.score.toFixed(2)}
                    </span>
                </>
            )

            objectElements.push(
                <div
                    style={{
                        display: "flex",
                        flexDirection: "row",
                        justifyContent: "space-between",
                    }}
                >
                    {inner}
                </div>,
            )
        })

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
                <td style={{ verticalAlign: "middle", width: "40px" }}>
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

    useEffect(() => {
        props.setLoading(eventsQuery.loading)
    }, [eventsQuery.loading, props])

    return (
        <Table striped bordered hover size="sm">
            <thead>
                <tr>
                    <th>Details</th>
                    <th>
                        Detections{" "}
                        <span style={{ fontWeight: "normal" }}>
                            {props.responsive && <br />}
                            (frames @ score)
                        </span>
                    </th>
                    <th colSpan={2}>Media</th>
                </tr>
            </thead>
            <tbody>{rows}</tbody>
        </Table>
    )
}
