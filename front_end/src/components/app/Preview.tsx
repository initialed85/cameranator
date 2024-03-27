import { Image } from "react-bootstrap"
import { fileHttpUrl } from "../../config"
import { Event } from "../../hasura/event"
import { Ref, useEffect, useState } from "react"
import { useInView } from "react-intersection-observer"

const ORIGINAL_WIDTH = 640
const ORIGINAL_HEIGHT = 360

const FOCUSED_WIDTH = ORIGINAL_WIDTH / 2
const FOCUSED_HEIGHT = ORIGINAL_HEIGHT / 2

const UNFOCUSED_WIDTH = FOCUSED_WIDTH / 2
const UNFOCUSED_HEIGHT = FOCUSED_HEIGHT / 2

export function adjustPath(path?: string): string {
    return fileHttpUrl + (path?.split("/srv/target_dir/")[1] || "")
}

export interface PreviewProps {
    event: Event
    focusEventUUID: string | null
}

export function Preview(props: PreviewProps) {
    const { ref, inView } = useInView({})
    const [alreadyLoaded, setAlreadyLoaded] = useState(false)

    useEffect(() => {
        if (inView && !alreadyLoaded) {
            setAlreadyLoaded(true)
        }
    }, [alreadyLoaded, inView])

    if (!props?.event?.thumbnail_image?.file_path) {
        return null
    }

    const focused = props.focusEventUUID === props.event.id

    return (
        <div
            id={props.event.id}
            ref={ref as Ref<HTMLDivElement>}
            style={{
                display: "flex",
                flexDirection: "row",
                justifyContent: "center",
                alignContent: "center",
                alignItems: "center",
                margin: 0,
                padding: 0,
                width: focused ? FOCUSED_WIDTH : UNFOCUSED_WIDTH,
                height: focused ? FOCUSED_HEIGHT : UNFOCUSED_HEIGHT,
            }}
        >
            {(inView || alreadyLoaded) && (
                <Image
                    id={props.event.id}
                    src={adjustPath(props.event.thumbnail_image.file_path)}
                    rounded
                    style={{
                        width: focused ? FOCUSED_WIDTH : UNFOCUSED_WIDTH,
                        height: focused ? FOCUSED_HEIGHT : UNFOCUSED_HEIGHT,
                    }}
                />
            )}
        </div>
    )
}
