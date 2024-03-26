import "bootstrap/dist/css/bootstrap.min.css"
import { Dispatch, SetStateAction } from "react"
import "./App.css"
import { Camera } from "../../hasura/camera"
import {
    Button,
    ButtonGroup,
    Dropdown,
    Nav,
    Navbar,
    ToggleButton,
    ToggleButtonGroup,
    FormControl,
    Spinner,
} from "react-bootstrap"
import moment from "moment"
import { Type } from "../../hasura/type"
import { fileHttpUrl } from "../../config"
import { Check } from "react-bootstrap-icons"

function getCameraButtons(
    cameras: Camera[] | null,
    camera: Camera | null,
    setCamera: Dispatch<SetStateAction<null>>,
    responsive: boolean,
) {
    if (!cameras?.length) {
        return null
    }

    const buttons: JSX.Element[] = []

    cameras.forEach((item) => {
        buttons.push(
            <ToggleButton
                size="sm"
                key={item.id}
                variant={"outline-secondary"}
                value={item.id}
                id={item.id}
                active={item.id === camera?.id}
                onClick={() => {
                    setCamera(item as any)
                }}
            >
                {item.name}
            </ToggleButton>,
        )
    })

    return (
        <ToggleButtonGroup
            name={"camera"}
            type={"radio"}
            style={{
                paddingLeft: responsive ? "5px" : "0px",
                paddingRight: responsive ? "5px" : "0px",
            }}
        >
            {buttons}
        </ToggleButtonGroup>
    )
}

// TODO
function getStreamButton(camera: Camera | null, responsive: boolean) {
    return (
        <ButtonGroup
            style={{
                paddingLeft: responsive ? "5px" : "0px",
                paddingRight: responsive ? "5px" : "0px",
            }}
        >
            <Button
                size="sm"
                variant={"outline-secondary"}
                disabled={!camera}
                href={
                    camera ? `${fileHttpUrl}stream/${camera.id}/stream/` : "#"
                }
                target={"_stream"}
            >
                Stream
            </Button>
        </ButtonGroup>
    )
}

function getDateDropdown(
    dates: moment.Moment[] | null,
    date: moment.Moment | null,
    setDate: Dispatch<SetStateAction<null>>,
    responsive: boolean,
) {
    if (!dates?.length) {
        return null
    }

    const eventDates: JSX.Element[] = []

    dates.forEach((item) => {
        const dateFriendly = item.format("YYYY-MM-DD")
        eventDates.push(
            <Dropdown.Item
                variant={"outline-secondary"}
                href={"#"}
                key={dateFriendly}
                id={dateFriendly}
                eventKey={dateFriendly}
                active={item.toISOString() === date?.toISOString()}
                onClick={() => {
                    setDate(item as any)
                }}
            >
                {dateFriendly}
            </Dropdown.Item>,
        )
    })

    return (
        <Dropdown
            style={{
                paddingLeft: "5px",
                paddingRight: responsive ? "5px" : "0px",
                width: "100%",
                marginTop: responsive ? "5px" : "0px",
            }}
        >
            <Dropdown.Toggle
                variant={"outline-secondary"}
                id="date"
                size="sm"
                style={{ width: "100%" }}
            >
                {date ? date.format("YYYY-MM-DD") : "Date"}
            </Dropdown.Toggle>

            <Dropdown.Menu variant={"outline-primary"}>
                {eventDates}
            </Dropdown.Menu>
        </Dropdown>
    )
}

function getObjectFilter(
    setObjectFilter: Dispatch<SetStateAction<string>>,
    responsive: boolean,
) {
    return (
        <FormControl
            id="objectFilter"
            type="text"
            size={"sm"}
            style={{
                marginLeft: "5px",
                width: responsive ? "98.05%" : "200px",
                marginTop: responsive ? "5px" : "0px",
                marginBottom: responsive ? "5px" : "0px",
                border: "1px solid #6c757d",
            }}
            onChange={(event) => {
                setObjectFilter(event.target.value)
            }}
        />
    )
}

export interface MenuProps {
    responsive: boolean
    cameras: Camera[] | null
    camera: Camera | null
    setCamera: Dispatch<SetStateAction<null>>
    dates: moment.Moment[] | null
    date: moment.Moment | null
    setDate: Dispatch<SetStateAction<null>>
    types: Type[] | null
    type: Type | null
    setType: Dispatch<SetStateAction<null>>
    setObjectFilter: Dispatch<SetStateAction<string>>
    isLoading: boolean
}

export function Menu(props: MenuProps) {
    return (
        <Navbar bg="light" expand="lg" style={{ fontSize: "10pt", padding: 0 }}>
            <Navbar.Brand
                href="#"
                style={{
                    fontSize: "14pt",
                    fontWeight: "bold",
                    marginLeft: "10px",
                    marginRight: "0px",
                    color: "gray",
                    width: "140px",
                }}
                onClick={() => {}}
            >
                Cameranator
                {props.isLoading ? (
                    <Spinner
                        size="sm"
                        animation="border"
                        style={{
                            color: "gray",
                            marginLeft: 5,
                        }}
                    />
                ) : (
                    <Check
                        style={{
                            width: "20px",
                            height: "20px",
                            marginLeft: "3px",
                        }}
                    />
                )}
            </Navbar.Brand>

            <Navbar.Toggle
                style={{
                    marginRight: "5px",
                    marginTop: props.responsive ? "5px" : "0px",
                    marginBottom: "5px",
                    color: "gray",
                }}
            />

            <Navbar.Collapse>
                <Nav>
                    {getCameraButtons(
                        props.cameras,
                        props.camera,
                        props.setCamera,
                        props.responsive,
                    )}
                    {/* {getStreamButton(props.camera)} */}
                    {getDateDropdown(
                        props.dates,
                        props.date,
                        props.setDate,
                        props.responsive,
                    )}
                </Nav>
                {getObjectFilter(props.setObjectFilter, props.responsive)}
            </Navbar.Collapse>
        </Navbar>
    )
}
