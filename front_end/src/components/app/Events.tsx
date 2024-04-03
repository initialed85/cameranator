import { Camera } from "../../hasura/camera"
import moment from "moment"
import { Type } from "../../hasura/type"
import { useQuery, useSubscription } from "@apollo/client"
import {
    Event,
    getEventsQuery,
    getEventsSubscription,
} from "../../hasura/event"
import { Modal, Table } from "react-bootstrap"
import { useEffect, useState } from "react"
import { adjustPath, Preview } from "./Preview"
import { CloudDownload, Play } from "react-bootstrap-icons"
import { Video } from "./Video"
import { Detection } from "../../hasura/detection"
import "./styles.css"
import { useInView } from "react-intersection-observer"

const MIN_SECONDS_SEEN = 2.0
const MIN_SCORE = 0.55

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
    const { ref, inView } = useInView({})
    const [eventById, setEventById] = useState<{ [n: number]: Event }>({})
    const [startTimestampLessThan, setStartTimestampLessThan] = useState<
        string | null
    >(null)

    const startTimestamp = props.date
        ? moment(`${props.date.local().format("YYYY-MM-DD")}T00:00:00+0800`)
        : null

    const endTimestamp = props.date
        ? moment(`${props.date.local().format("YYYY-MM-DD")}T23:59:59+0800`)
        : null

    const eventsSubscription = useSubscription(getEventsSubscription(), {
        variables: {
            startTimestamp: startTimestamp,
            endTimestamp: endTimestamp,
            sourceCameraId: props.camera?.id,
        },
    })

    const eventsQuery = useQuery(getEventsQuery(), {
        variables: {
            camera: props.camera,
            date: props.date,
            type: props.type,
            startTimestampLessThan: startTimestampLessThan,
        },
    })

    useEffect(() => {
        const thisEventById: { [n: number]: Event } = { ...eventById }

        if (
            !eventsSubscription.loading &&
            eventsSubscription.data?.event?.length
        ) {
            eventsSubscription.data?.event?.forEach((event: Event) => {
                thisEventById[parseInt(event.id)] = event
            })
        }

        if (!eventsQuery.loading && eventsQuery.data?.event?.length) {
            eventsQuery.data?.event?.forEach((event: Event) => {
                thisEventById[parseInt(event.id)] = event
            })
        }

        const thisEventByIdHash = JSON.stringify(
            Object.keys(thisEventById)
                .map((k) => parseInt(k))
                .sort(),
        )

        const eventByIdHash = JSON.stringify(
            Object.keys(eventById)
                .map((k) => parseInt(k))
                .sort(),
        )

        if (thisEventByIdHash !== eventByIdHash) {
            setEventById({ ...eventById, ...thisEventById })
        }
    }, [
        eventById,
        eventsQuery.data?.event,
        eventsQuery.loading,
        eventsSubscription.data?.event,
        eventsSubscription.loading,
        startTimestampLessThan,
        inView,
        props,
    ])

    useEffect(() => {
        props.setLoading(eventsSubscription.loading || eventsQuery.loading)
    }, [eventsQuery.loading, eventsSubscription.loading, props])

    const [focusEventUUID, setFocusEventUUID] = useState(null)

    const [showModal, setShowModal] = useState(false)
    const [event, setEvent] = useState<Event | null>(null)

    const rawObjectFilter = (props.objectFilter.trim() || "")
        .replaceAll(" ", "")
        .toLocaleLowerCase()

    const objectFilterParts = rawObjectFilter
        .split(",")
        .filter((x) => x.trim() !== "")

    const matchesObjectFilter = (
        className: string,
        isOutlier: boolean = false,
    ): boolean => {
        if (!objectFilterParts.length) {
            return true
        }

        const matches = objectFilterParts.map((x: string): boolean => {
            const excludeOutliers = x.includes("!")

            if (excludeOutliers && isOutlier) {
                return false
            }

            if (className.includes(x.replaceAll("!", ""))) {
                return true
            }

            return false
        })

        return matches.some((x) => !!x)
    }

    const rows: JSX.Element[] = []

    const events: Event[] = Object.values(eventById)
        .filter((event: Event) => {
            if (event.source_camera.id !== props.camera?.id) {
                return false
            }

            const thisStartTimestamp = moment(event.start_timestamp)
            if (thisStartTimestamp.isBefore(startTimestamp)) {
                return false
            }

            const thisEndTimestamp = moment(event.end_timestamp)
            if (thisEndTimestamp.isAfter(endTimestamp)) {
                return false
            }

            return true
        })
        .sort((a: Event, b: Event) => {
            if (a.start_timestamp < b.start_timestamp) {
                return 1
            } else if (a.start_timestamp > b.start_timestamp) {
                return -1
            } else {
                return 0
            }
        })

    if (inView && !eventsSubscription.loading && !eventsQuery.loading) {
        const thisStartTimestampLessThan = events
            .map((event: Event) => event.start_timestamp)
            .sort()[0]

        if (
            thisStartTimestampLessThan &&
            thisStartTimestampLessThan !== startTimestampLessThan
        ) {
            setStartTimestampLessThan(thisStartTimestampLessThan)
        }
    }

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

            if (objectFilterParts.length) {
                // 20 fps / 4 stride frames = seconds seen
                const outlier =
                    detection.count / (20 / 4) < MIN_SECONDS_SEEN ||
                    detection.score < MIN_SCORE

                if (matchesObjectFilter(detection.class_name, outlier)) {
                    filteredDetectionByClassName.set(
                        detection.class_name,
                        detection,
                    )
                }
            }
        })

        if (
            objectFilterParts.length &&
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
                detection.score < MIN_SCORE

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

        detectionSummaries.forEach((detectionSummary, i: number) => {
            let color = "black"
            let textDecoration = "none"
            if (detectionSummary.outlier) {
                color = "gray"
                textDecoration = "line-through"
            }

            if (
                objectFilterParts.length &&
                matchesObjectFilter(
                    detectionSummary.className,
                    detectionSummary.outlier,
                )
            ) {
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
                    key={`${event.id}-${i}`}
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
                        key={event.id}
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
                        key={event.id}
                        style={{
                            display: "flex",
                            flexDirection: "column",
                            paddingLeft: "5px",
                            paddingRight: "5px",
                            overflow: "scroll",
                            height:
                                focusEventUUID === event.id ? "185px" : "95px",
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
                        key={event.id}
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
                    <div
                        key={event.id}
                        style={{
                            display: "flex",
                            flexDirection: "column",
                            alignItems: "center",
                            justifyContent: "space-between",
                            height: "100%",
                        }}
                    >
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
                        <Play
                            style={{
                                width: "25px",
                                height: "25px",
                                color: "gray",
                            }}
                            role={"button"}
                            onClick={() => {
                                setEvent(event)
                                setShowModal(true)
                            }}
                        />
                    </div>
                </td>
            </tr>,
        )
    })

    return (
        <>
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

            <span
                ref={ref}
                style={{
                    color: "#aaaaaa",
                    textAlign: "center",
                    fontSize: "6pt",
                }}
            />

            <Modal
                contentClassName={"videoModal"}
                dialogClassName={"videoModal"}
                backdropClassName={"videoModal"}
                show={showModal}
                onHide={() => {
                    setShowModal(false)
                }}
                size={"xl"}
            >
                <Modal.Body
                    style={{
                        padding: 2,
                        margin: 0,
                    }}
                >
                    {event && <Video event={event} />}
                </Modal.Body>
            </Modal>
        </>
    )
}
