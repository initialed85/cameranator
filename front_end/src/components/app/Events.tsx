import { Camera } from "../../hasura/camera";
import moment from "moment";
import { Type } from "../../hasura/type";
import { useSubscription } from "@apollo/client";
import { Event, getEventsQuery } from "../../hasura/event";
import { Table } from "react-bootstrap";
import { useState } from "react";
import { adjustPath, Preview } from "./Preview";

export interface EventsProps {
  camera: Camera | null;
  date: moment.Moment | null;
  type: Type | null;
  responsive: boolean;
}

export function Events(props: EventsProps) {
  const eventsQuery = useSubscription(
    getEventsQuery(props.camera!, props.date!, props.type!)
  );

  const [focusEventUUID, setFocusEventUUID] = useState(null);

  // TODO
  // eslint-disable-next-line
  const [confirmedFocusEventUUID, setConfirmedFocusEventUUID] = useState(null);

  const rows: JSX.Element[] = [];

  (eventsQuery?.data?.event || []).forEach((event: Event) => {
    const startTimestamp = moment.utc(event.start_timestamp);
    const endTimestamp = moment.utc(event.end_timestamp);

    const duration = moment.duration(endTimestamp.diff(event.start_timestamp));
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
        <td
          width={!props.responsive ? 330 : 170}
          onMouseOver={(e) => {
            if (props.responsive) {
              return;
            }

            setFocusEventUUID(event.uuid as any);

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
              return;
            }

            setFocusEventUUID(null);

            // TODO
            // setConfirmedFocusEventUUID(null);
          }}
          onClick={(e) => {
            if (props.responsive) {
              return;
            }

            setFocusEventUUID(null);

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
            High-res
          </a>
          <br />
          <a
            target={`_low_quality_video_${event.uuid}`}
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
          <th>Period</th>
          <th>Duration</th>
          <th>Image</th>
          <th>Download</th>
        </tr>
      </thead>
      <tbody>{rows}</tbody>
    </Table>
  );
}
