import { Camera } from "../../hasura/camera"
import moment from "moment/moment"
import { Type } from "../../hasura/type"
import { Events } from "./Events"

export interface ContentProps {
    responsive: boolean
    camera: Camera | null
    date: moment.Moment | null
    type: Type | null
    objectFilter: string
    setLoading: (x: boolean) => void
}

export function Content(props: ContentProps) {
    if (!(props.camera && props.date && props.type)) {
        return null
    }

    if (props.type.is_stream) {
        return null
    }

    return (
        <Events
            responsive={props.responsive}
            camera={props.camera}
            date={props.date}
            type={props.type}
            objectFilter={props.objectFilter}
            setLoading={props.setLoading}
        />
    )
}
