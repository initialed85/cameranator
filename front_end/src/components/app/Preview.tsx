import { Image } from "react-bootstrap"
import { fileHttpUrl } from "../../config"
import { Event } from "../../hasura/event"
import { Ref, useEffect, useState } from "react"
import { useInView } from "react-intersection-observer"

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
                width: focused ? 320 : 120,
                height: focused ? 180 : 90,
            }}
        >
            {(inView || alreadyLoaded) && (
                <Image
                    id={props.event.id}
                    src={adjustPath(props.event.thumbnail_image.file_path)}
                    rounded
                    style={{
                        width: focused ? 320 : 120,
                        height: focused ? 180 : 90,
                    }}
                />
            )}
        </div>
    )
}
