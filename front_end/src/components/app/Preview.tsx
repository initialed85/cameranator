import { Image } from "react-bootstrap";
import { fileHttpUrl } from "../../config";
import { Event } from "../../hasura/event";

export function adjustPath(path: string): string {
  return fileHttpUrl + path.split("/srv/target_dir/")[1];
}

export interface PreviewProps {
  event: Event;
  focusEventUUID: string | null;
  confirmedFocusEventUUID: string | null;
}

export function Preview(props: PreviewProps) {
  // TODO
  // if (props?.confirmedFocusEventUUID === props.event.uuid) {
  //   return (
  //     <video width={640} height={480} autoPlay={true}>
  //       <source src={adjustPath(props.event.low_quality_video.file_path)} />
  //     </video>
  //   );
  // }

  return (
    <Image
      src={adjustPath(props.event.low_quality_image.file_path)}
      rounded
      style={{ width: props.focusEventUUID === props.event.uuid ? 320 : 160 }}
    />
  );
}
