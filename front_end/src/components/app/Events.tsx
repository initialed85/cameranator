import { Camera } from "../../hasura/camera"
import moment from "moment"
import { Type } from "../../hasura/type"
import { useSubscription } from "@apollo/client"
import { Event, getEventsQuery, Object } from "../../hasura/event"
import { Table } from "react-bootstrap"
import { useState } from "react"
import { adjustPath, Preview } from "./Preview"

export interface EventsProps {
    responsive: boolean
    camera: Camera | null
    date: moment.Moment | null
    type: Type | null
    objectFilter: string
}

interface ObjectWithContext {
    object: Object
    count: number
    durationMilliseconds: number
}

export function Events(props: EventsProps) {
    const eventsQuery = useSubscription(
        getEventsQuery(props.camera!, props.date!, props.type!),
    )

    const [focusEventUUID, setFocusEventUUID] = useState(null)

    // TODO
    // eslint-disable-next-line
    const [confirmedFocusEventUUID, setConfirmedFocusEventUUID] = useState(null)

    const rows: JSX.Element[] = []

    ;(eventsQuery?.data?.event || []).forEach((event: Event) => {
        const startTimestamp = moment.utc(event.start_timestamp)
        const endTimestamp = moment.utc(event.end_timestamp)

        const duration = moment.duration(
            endTimestamp.diff(event.start_timestamp),
        )

        const objectFilter = (props.objectFilter.trim() || "").toLowerCase()

        const objectsByHash: Map<string, Object[]> = new Map()
        event.objects.forEach((object: Object) => {
            let sanitised = {
                ...object,
            }
            sanitised.tracked_object_id = 0
            sanitised.start_timestamp = ""
            sanitised.end_timestamp = ""

            const hash = JSON.stringify(sanitised)

            const objects: Object[] = objectsByHash.get(hash) || []
            objects.push(object)
            objectsByHash.set(hash, objects)
        })

        let objectsWithContext: ObjectWithContext[] = []
        objectsByHash.forEach((objects: Object[], hash: string) => {
            let count = 0
            let durationMilliseconds = 0
            objects.forEach((object: Object) => {
                count++
                durationMilliseconds +=
                    new Date(object.end_timestamp).valueOf() -
                    new Date(object.start_timestamp).valueOf()
            })

            objectsWithContext.push({
                object: JSON.parse(hash),
                count,
                durationMilliseconds,
            } as ObjectWithContext)
        })

        objectsWithContext = objectsWithContext.sort(
            (a: ObjectWithContext, b: ObjectWithContext) => {
                if (a.count < b.count) {
                    return 0
                } else if (a.count > b.count) {
                    return -1
                } else {
                    return 0
                }
            },
        )

        const objectElements: JSX.Element[] = []
        objectsWithContext.forEach((objectWithContext, i) => {
            objectElements.push(
                <span key={i}>
                    {objectWithContext.object.detected_class_name} (
                    {objectWithContext.object.detected_class_id}) x{" "}
                    {objectWithContext.count}
                </span>,
            )
        })

        const filteredObjectsWithContext = objectsWithContext.filter(
            (objectWithContext: ObjectWithContext) => {
                if (objectFilter.length) {
                    if (
                        !objectWithContext.object.detected_class_name
                            .toLowerCase()
                            .includes(objectFilter)
                    ) {
                        return false
                    }
                }

                return true
            },
        )

        if (objectFilter.length && !filteredObjectsWithContext.length) {
            return
        }

        rows.push(
            <tr key={event.uuid}>
                <td>
                    {startTimestamp.local().format("HH:mm:ss")}
                    {props.responsive ? <br /> : " to "}
                    {endTimestamp.local().format("HH:mm:ss")}
                </td>
                <td>
                    {duration.minutes()}m{duration.seconds()}s
                </td>
                <td>
                    <div style={{ display: "flex", flexDirection: "column" }}>
                        {objectElements}
                    </div>
                </td>
                <td
                    width={!props.responsive ? 330 : 170}
                    onMouseOver={(e) => {
                        if (props.responsive) {
                            return
                        }

                        setFocusEventUUID(event.uuid as any)

                        // TODO
                        // setTimeout(() => {
                        //   if (focusEventUUID !== event.uuid) {
                        //     return;
                        //   }
                        //   setConfirmedFocusEventUUID(event.uuid as any);
                        // }, 1_000);
                    }}
                    onMouseOut={(e) => {
                        if (props.responsive) {
                            return
                        }

                        setFocusEventUUID(null)

                        // TODO
                        // setConfirmedFocusEventUUID(null);
                    }}
                    onClick={(e) => {
                        if (props.responsive) {
                            return
                        }

                        setFocusEventUUID(null)

                        // TODO
                        // setConfirmedFocusEventUUID(null);
                    }}
                >
                    <a
                        target={`_high_quality_image_${event.uuid}`}
                        rel={"noreferrer"}
                        href={adjustPath(event.high_quality_image.file_path)}
                    >
                        <Preview
                            event={event}
                            focusEventUUID={focusEventUUID}
                            confirmedFocusEventUUID={confirmedFocusEventUUID}
                        />
                    </a>
                </td>
                <td>
                    <a
                        target={`_high_quality_video_${event.uuid}`}
                        rel={"noreferrer"}
                        href={adjustPath(event.high_quality_video.file_path)}
                    >
                        Original
                    </a>
                    <br />
                    <a
                        target={`_low_quality_video_${event.uuid}`}
                        rel={"noreferrer"}
                        href={adjustPath(
                            event.high_quality_video.file_path.replaceAll(
                                ".mp4",
                                "_out.mp4",
                            ),
                        )}
                    >
                        Tracked objects
                    </a>
                </td>
            </tr>,
        )
    })

    return (
        <Table striped bordered hover size="sm">
            <thead>
                <tr>
                    <th>Period</th>
                    <th>Duration</th>
                    <th>Tracked objects</th>
                    <th>Image</th>
                    <th>Download</th>
                </tr>
            </thead>
            <tbody>{rows}</tbody>
        </Table>
    )
}
