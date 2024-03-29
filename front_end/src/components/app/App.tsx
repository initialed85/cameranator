import "bootstrap/dist/css/bootstrap.min.css"
import { useCallback, useEffect, useState } from "react"
import "./App.css"
import { CAMERAS, Camera } from "../../hasura/camera"
import { useSubscription } from "@apollo/client"
import { Container, Row } from "react-bootstrap"
import { EVENT_DATES, getDeduplicatedDates } from "../../hasura/event_date"
import { Menu } from "./Menu"
import { TYPES, Type } from "../../hasura/type"
import { Content } from "./Content"
import { useDebouncedCallback } from "use-debounce"
import { Moment } from "moment"

function App() {
    const camerasQuery = useSubscription(CAMERAS)
    const datesQuery = useSubscription(EVENT_DATES)
    const types = TYPES

    const deduplicatedDates = getDeduplicatedDates(datesQuery?.data)

    const [responsive, setResponsive] = useState(window.innerWidth < 992)
    const [camera, setCamera] = useState<Camera | null>(null)
    const [date, setDate] = useState<Moment | null>(null)
    const [type, setType] = useState<Type | null>(null)
    const [objectFilter, setObjectFilter] = useState("")
    const debouncedSetObjectFilter = useDebouncedCallback(
        (x) => setObjectFilter(x),
        100,
    )

    const [isLoading, setIsLoading] = useState(false)

    const setLoading = useCallback(
        (loading: boolean) => {
            if (loading === isLoading) {
                return
            }

            setIsLoading(loading)
        },
        [isLoading],
    )

    useEffect(() => {
        const handleResize = () => {
            setResponsive(window.innerWidth < 992)
        }

        window.addEventListener("resize", () => {
            handleResize()
        })

        if (!type && types?.length) {
            setType(types[0] as any)
        }
        if (!camera && camerasQuery?.data?.camera) {
            setCamera(camerasQuery?.data?.camera[0])
        }

        if (!date && deduplicatedDates?.length) {
            setDate(deduplicatedDates[0] as any)
        }

        if (!type && types?.length) {
            setType(types[0] as any)
        }

        return () => {
            window.removeEventListener("resize", handleResize)
        }
    }, [
        camera,
        camerasQuery?.data?.camera,
        date,
        deduplicatedDates,
        type,
        types,
    ])

    return (
        <Container>
            <Row>
                <Menu
                    responsive={responsive}
                    cameras={camerasQuery?.data?.camera}
                    camera={camera}
                    setCamera={(x) => setCamera(x)}
                    dates={deduplicatedDates}
                    date={date}
                    setDate={(x) => setDate(x)}
                    types={types}
                    type={type}
                    setType={(x) => setType(x)}
                    setObjectFilter={(x) => debouncedSetObjectFilter(x)}
                    isLoading={
                        camerasQuery?.loading ||
                        datesQuery?.loading ||
                        isLoading
                    }
                />
            </Row>
            <Row>
                <Content
                    responsive={responsive}
                    camera={camera}
                    date={date}
                    type={type}
                    objectFilter={objectFilter}
                    setLoading={(x) => setLoading(x)}
                />
            </Row>
        </Container>
    )
}

export default App
