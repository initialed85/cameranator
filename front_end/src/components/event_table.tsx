import { Event } from "../persistence/event";
import React from "react";
import { Image, Table } from "react-bootstrap";
import moment from "moment";

// TODO: make this configurable
const urlPrefix = "/";

function adjustPath(path: string): string {
    return urlPrefix + path.split("/srv/target_dir/")[1];
}

interface EventTableProps {
    events: Event[];
}

export class EventTable extends React.Component<EventTableProps, any> {
    getEvents(): any {
        let rows: any[] = [];

        this.props.events.forEach((event) => {
            const duration = moment.duration(
                event.end_timestamp.diff(event.start_timestamp)
            );
            rows.push(
                <tr key={event.uuid}>
                    <td>{event.start_timestamp.local().format("HH:mm:ss")}</td>
                    <td>{event.end_timestamp.local().format("HH:mm:ss")}</td>
                    <td>
                        {duration.minutes()}:
                        {duration.seconds().toString().padStart(2, "0")}
                    </td>
                    <td>{event.source_camera.name}</td>
                    <td style={{ width: 320 }}>
                        <a
                            target={"_blank"}
                            href={adjustPath(
                                event.high_quality_image.file_path
                            )}
                        >
                            <Image
                                src={adjustPath(
                                    event.low_quality_image.file_path
                                )}
                                rounded
                                style={{ width: "100%" }}
                            />
                        </a>
                    </td>
                    <td>
                        <a
                            target="_blank"
                            href={adjustPath(
                                event.high_quality_video.file_path
                            )}
                        >
                            High-res
                        </a>
                        <br />
                        <a
                            target="_blank"
                            href={adjustPath(event.low_quality_video.file_path)}
                        >
                            Low-res
                        </a>
                    </td>
                </tr>
            );
        });

        return <tbody>{rows}</tbody>;
    }

    render() {
        return (
            <Table striped bordered hover size="sm">
                <thead>
                    <tr>
                        <th>Start</th>
                        <th>End</th>
                        <th>Duration</th>
                        <th>Camera</th>
                        <th>Image</th>
                        <th>Download</th>
                    </tr>
                </thead>
                {this.getEvents()}
            </Table>
        );
    }
}
