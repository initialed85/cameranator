import VideoJS from "./VideoJS"
import { Event } from "../../hasura/event"
import { adjustPath } from "./Preview"
import { useQuery } from "@apollo/client"
import {
    Detection,
    Point,
    getDetectionsForEventQuery,
} from "../../hasura/detection"
import { useEffect, useRef } from "react"

export interface VideoProps {
    event: Event
}

export function Video(props: VideoProps) {
    const { loading, error, data } = useQuery(
        getDetectionsForEventQuery(props.event.id),
    )

    const enrichedDetectionsRef = useRef<Detection[]>([])
    const canvasRef = useRef<HTMLCanvasElement | null>(null)
    const readyRef = useRef(false)

    useEffect(() => {
        if (enrichedDetectionsRef.current?.length) {
            return
        }

        let enrichedDetections: Detection[] = []

        const detections: Detection[] = data?.detections || []
        detections.forEach((detection: Detection) => {
            const rawBoundingBox = detection.bounding_box as string
            const boundingBoxPoints: Point[] = []
            rawBoundingBox
                .replaceAll("((", "")
                .replaceAll("))", "")
                .split("),(")
                .forEach((part: string, i: number) => {
                    const [rawX, rawY] = part.split(",")

                    const point = {
                        x: parseFloat(rawX),
                        y: parseFloat(rawY),
                    } as Point

                    boundingBoxPoints.push(point)
                })

            const rawCentroid = detection.centroid as string
            const [rawX, rawY] = rawCentroid
                .replaceAll("(", "")
                .replaceAll(")", "")
                .split(",")
            const centroidPoint = {
                x: parseFloat(rawX),
                y: parseFloat(rawY),
            } as Point

            const enrichedDetection = {
                ...detection,
                timestampMilliseconds: Date.parse(detection.timestamp),
                boundingBoxPoints,
                centroidPoint,
            }

            enrichedDetections.push(enrichedDetection)
        })

        enrichedDetectionsRef.current = enrichedDetections
    }, [data?.detections])

    const absoluteTimeMillisecondsRef = Date.parse(props.event.start_timestamp)

    if (error) {
        console.warn(error)
        return (
            <div style={{ fontWeight: "bold", color: "red" }}>
                ERROR: {JSON.stringify(error)}
            </div>
        )
    }

    if (loading) {
        return null
    }

    return (
        <>
            <canvas
                ref={canvasRef}
                style={{
                    position: "fixed",
                    display: "block",
                    zIndex: 500,
                    left: 0,
                    top: 0,
                    width: 0,
                    height: 0,
                    cursor: "not-allowed",
                    pointerEvents: "none",
                }}
                width={1920}
                height={1080}
            />
            <VideoJS
                options={{
                    autoplay: true,
                    controls: true,
                    responsive: true,
                    fluid: true,
                    ratio: "16:9",
                    inactivityTimeout: 0,
                    playsinline: true,
                    preload: "auto",
                    enableSmoothSeeking: true,
                    sources: [
                        {
                            src: adjustPath(
                                props.event.original_video.file_path,
                            ),
                            type: "video/mp4",
                        },
                    ],
                }}
                onReady={() => {
                    readyRef.current = true
                }}
                onTimeUpdate={(
                    left: number,
                    top: number,
                    width: number,
                    height: number,
                    relativeTimeMilliseconds: number,
                ) => {
                    if (!canvasRef.current) {
                        return
                    }

                    const canvas = canvasRef.current as HTMLCanvasElement

                    const firstUpdate = canvas.style.left === "0px"

                    canvas.style.left = `${left}px`
                    canvas.style.top = `${top}px`
                    canvas.style.width = `${width}px`
                    canvas.style.height = `${height}px`

                    if (!readyRef?.current) {
                        return
                    }

                    const ctx = canvas.getContext(
                        "2d",
                    ) as CanvasRenderingContext2D
                    if (!ctx) {
                        return
                    }

                    ctx.clearRect(0, 0, 1920, 1080)

                    if (firstUpdate) {
                        return
                    }

                    const absoluteTimeMilliseconds =
                        absoluteTimeMillisecondsRef + relativeTimeMilliseconds

                    const enrichedDetections =
                        enrichedDetectionsRef.current || []
                    enrichedDetections.forEach((detection: Detection) => {
                        const deltaMilliseconds =
                            absoluteTimeMilliseconds -
                            detection.timestampMilliseconds

                        if (
                            deltaMilliseconds < 0 ||
                            deltaMilliseconds > 5_000
                        ) {
                            return
                        }

                        const topLeft = detection.boundingBoxPoints[0]
                        const bottomRight = detection.boundingBoxPoints[2]

                        ctx.lineWidth = 2
                        ctx.strokeStyle = `rgba(255, 0, 0, 0.75)`

                        ctx.fillStyle = `rgba(255, 255, 255, 0.75)`
                        ctx.font = "18px sans-serif"
                        ctx.textAlign = "left"

                        if (deltaMilliseconds <= 200) {
                            ctx.fillText(
                                `${
                                    detection.class_name
                                } @ ${detection.score.toFixed(3)}`,
                                topLeft.x + 3,
                                bottomRight.y - 4,
                            )

                            ctx.strokeRect(
                                topLeft.x,
                                topLeft.y,
                                Math.abs(bottomRight.x - topLeft.x),
                                Math.abs(bottomRight.y - topLeft.y),
                            )
                        }

                        const color = (1.0 - deltaMilliseconds / 5_000) * 255
                        const alpha = 1.0 - deltaMilliseconds / 5_000

                        ctx.strokeStyle = `rgba(${color}, ${color}, ${color}, ${alpha})`

                        ctx.beginPath()

                        ctx.arc(
                            detection.centroidPoint.x,
                            detection.centroidPoint.y,
                            5,
                            0,
                            Math.PI * 2,
                        )

                        ctx.stroke()
                    })
                }}
            />
        </>
    )
}
