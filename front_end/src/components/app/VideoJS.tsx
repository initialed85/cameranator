import { MutableRefObject, useEffect, useRef } from "react"
import videojs from "video.js"
import "video.js/dist/video-js.css"

interface VideoJSProps {
    options: any
    onReady: () => void
    onTimeUpdate: (
        left: number,
        top: number,
        width: number,
        height: number,
        relativeTimeMilliseconds: number,
    ) => void
}

export const VideoJS = (props: VideoJSProps) => {
    const videoRef: MutableRefObject<any> = useRef(null)
    const playerRef: MutableRefObject<any> = useRef(null)

    const { options, onReady, onTimeUpdate } = props

    useEffect(() => {
        if (!playerRef.current) {
            const videoElement = document.createElement("video-js")

            videoElement.classList.add("vjs-big-play-centered")
            videoElement.classList.add("vjs-has-started")
            videoElement.style.border = "0"
            videoElement.style.zIndex = "300"

            videoRef.current.appendChild(videoElement)

            const player = (playerRef.current = videojs(
                videoElement,
                options,
                () => {
                    const handleRequestAnimationFrame = () => {
                        player.controls(true)
                        player.usingNativeControls(false)

                        const relativeTimeSeconds = player.currentTime()
                        if (relativeTimeSeconds === undefined) {
                            return
                        }

                        const rect = videoElement.getBoundingClientRect()

                        onTimeUpdate(
                            rect.left,
                            rect.top,
                            rect.width,
                            rect.height,
                            relativeTimeSeconds * 1_000.0,
                        )

                        player.requestAnimationFrame(() => {
                            handleRequestAnimationFrame()
                        })
                    }

                    setTimeout(() => {
                        player.requestAnimationFrame(() => {
                            handleRequestAnimationFrame()
                        })
                    }, 3_00)

                    if (!onReady) {
                        return
                    }

                    onReady()
                },
            ))
        } else {
            const player = playerRef.current

            player.autoplay(options.autoplay)
            player.src(options.sources)
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    useEffect(() => {
        const player = playerRef.current

        return () => {
            if (player && !player.isDisposed()) {
                player.dispose()
                playerRef.current = null
            }
        }
    }, [playerRef])

    return (
        <div
            data-vjs-player
            style={{ margin: 0, padding: 0, border: 0, zIndex: 100 }}
        >
            <div
                ref={videoRef}
                style={{ margin: 0, padding: 0, border: 0, zIndex: 200 }}
            />
        </div>
    )
}

export default VideoJS
