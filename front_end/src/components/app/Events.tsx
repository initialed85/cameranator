import { Camera } from "../../hasura/camera";
import moment from "moment";
import { Type } from "../../hasura/type";
import { useQuery } from "@apollo/client";
import { Event, getEventsQuery } from "../../hasura/event";
import { Table } from "react-bootstrap";
import { useState } from "react";
import { adjustPath, Preview } from "./Preview";

export interface EventsProps {
  camera: Camera | null;
  date: moment.Moment | null;
  type: Type | null;
}

export function Events(props: EventsProps) {
  const eventsQuery = useQuery(
    getEventsQuery(props.camera!, props.date!, props.type!)
  );

  const [focusEventUUID, setFocusEventUUID] = useState(null);
  const [confirmedFocusEventUUID, setConfirmedFocusEventUUID] = useState(null);

  const rows: JSX.Element[] = [];

  (eventsQuery?.data?.event || []).forEach((event: Event) => {
    const startTimestamp = moment.utc(event.start_timestamp);
    const endTimestamp = moment.utc(event.end_timestamp);

    const duration = moment.duration(endTimestamp.diff(event.start_timestamp));
    rows.push(
      <tr key={event.uuid}>
        <td>{startTimestamp.local().format("HH:mm:ss")}</td>
        <td>{endTimestamp.local().format("HH:mm:ss")}</td>
        <td>
          {duration.minutes()}m{duration.seconds()}s
        </td>
        {/*<td>{event.source_camera.name}</td>*/}
        <td
          style={{
            width: 160,
          }}
          onMouseOver={(e) => {
            setFocusEventUUID(event.uuid as any);

            // TODO: disabled until unwanted margin fixed
            // setTimeout(() => {
            //   if (focusEventUUID !== event.uuid) {
            //     return;
            //   }
            //   setConfirmedFocusEventUUID(event.uuid as any);
            // }, 1_000);
          }}
          onMouseOut={(e) => {
            setFocusEventUUID(null);

            // TODO: disabled until unwanted margin fixed
            // setConfirmedFocusEventUUID(null);
          }}
          onClick={(e) => {
            setFocusEventUUID(null);

            // TODO: disabled until unwanted margin fixed
            // setConfirmedFocusEventUUID(null);
          }}
        >
          <a
            target={"_blank"}
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
            target="_blank"
            rel={"noreferrer"}
            href={adjustPath(event.high_quality_video.file_path)}
          >
            High-res
          </a>
          <br />
          <a
            target="_blank"
            rel={"noreferrer"}
            href={adjustPath(event.low_quality_video.file_path)}
          >
            Low-res
          </a>
        </td>
      </tr>
    );
  });

  return (
    <Table striped bordered hover size="sm">
      <thead>
        <tr>
          <th>Start</th>
          <th>End</th>
          <th>Duration</th>
          {/*<th>Camera</th>*/}
          <th>Image</th>
          <th>Download</th>
        </tr>
      </thead>
      <tbody>{rows}</tbody>
    </Table>
  );
}
