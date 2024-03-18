import { Image } from "react-bootstrap"
import { fileHttpUrl } from "../../config"
import { Event } from "../../hasura/event"

export function adjustPath(path?: string): string {
    return fileHttpUrl + (path?.split("/srv/target_dir/")[1] || "")
}

export interface PreviewProps {
    event: Event
    focusEventUUID: string | null
}

export function Preview(props: PreviewProps) {
    if (!props?.event?.thumbnail_image?.file_path) {
        return null
    }

    return (
        <Image
            src={adjustPath(props.event.thumbnail_image.file_path)}
            rounded
            style={{
                width: props.focusEventUUID === props.event.id ? 320 : 120,
            }}
        />
    )
}
