import { Camera } from "../../hasura/camera";
import moment from "moment/moment";
import { Type } from "../../hasura/type";
import { Events } from "./Events";

export interface ContentProps {
  camera: Camera | null;
  date: moment.Moment | null;
  type: Type | null;
}

export function Content(props: ContentProps) {
  if (!(props.camera && props.date && props.type)) {
    return null;
  }

  if (props.type.is_stream) {
    return null;
  }

  return <Events camera={props.camera} date={props.date} type={props.type} />;
}
