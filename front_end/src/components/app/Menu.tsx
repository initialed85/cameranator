import "bootstrap/dist/css/bootstrap.min.css"
import { Dispatch, SetStateAction } from "react"
import "./App.css"
import { Camera } from "../../hasura/camera"
import {
    Dropdown,
    Nav,
    Navbar,
    ToggleButton,
    ToggleButtonGroup,
    FormControl,
    Spinner,
    Button,
} from "react-bootstrap"
import moment, { Moment } from "moment"
import { Type } from "../../hasura/type"
import { Check } from "react-bootstrap-icons"
import { fileHttpUrl } from "../../config"

function getCameraButtons(
    cameras: Camera[] | null,
    camera: Camera | null,
    setCamera: Dispatch<SetStateAction<Camera | null>>,
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

function getDateDropdown(
    dates: moment.Moment[] | null,
    date: moment.Moment | null,
    setDate: Dispatch<SetStateAction<Moment | null>>,
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

function getStreamButton(camera: Camera | null, responsive: boolean) {
    return (
        <div
            style={{
                paddingLeft: "5px",
                paddingRight: "0px",
                marginTop: responsive ? "5px" : "0px",
                width: responsive ? "100%" : "150%",
            }}
        >
            <Button
                size="sm"
                variant={"outline-secondary"}
                style={{ width: "100%" }}
                disabled={!!!camera?.name}
                href={fileHttpUrl + "streams/" + camera?.name}
                target={camera?.name}
            >
                Stream
            </Button>
        </div>
    )
}

function getObjectFilter(
    setObjectFilter: Dispatch<SetStateAction<string>>,
    responsive: boolean,
) {
    return (
        <div
            style={{
                paddingLeft: "5px",
                paddingRight: "5px",
                width: responsive ? "100%" : "150%",
            }}
        >
            <FormControl
                id="objectFilter"
                type="text"
                size={"sm"}
                style={{
                    marginTop: responsive ? "5px" : "0px",
                    marginBottom: responsive ? "5px" : "0px",
                    border: "1px solid #6c757d",
                }}
                onChange={(event) => {
                    setObjectFilter(event.target.value)
                }}
            />
        </div>
    )
}

export interface MenuProps {
    responsive: boolean
    cameras: Camera[] | null
    camera: Camera | null
    setCamera: Dispatch<SetStateAction<Camera | null>>
    dates: moment.Moment[] | null
    date: moment.Moment | null
    setDate: Dispatch<SetStateAction<Moment | null>>
    types: Type[] | null
    type: Type | null
    setType: Dispatch<SetStateAction<Type | null>>
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
                    <div style={{ display: "flex", flexDirection: "row" }}>
                        {getStreamButton(props.camera, props.responsive)}
                        {getDateDropdown(
                            props.dates,
                            props.date,
                            props.setDate,
                            props.responsive,
                        )}
                    </div>
                    {getObjectFilter(props.setObjectFilter, props.responsive)}
                </Nav>
            </Navbar.Collapse>
        </Navbar>
    )
}
