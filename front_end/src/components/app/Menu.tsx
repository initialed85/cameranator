import "bootstrap/dist/css/bootstrap.min.css"
import React, { Dispatch, SetStateAction } from "react"
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
} from "react-bootstrap"
import moment from "moment"
import { Type } from "../../hasura/type"
import { fileHttpUrl } from "../../config"

function getTypeButtons(
    types: Type[] | null,
    type: Type | null,
    setType: Dispatch<SetStateAction<null>>,
) {
    if (!types?.length) {
        return null
    }

    const buttons: JSX.Element[] = []

    types.forEach((item) => {
        buttons.push(
            <ToggleButton
                size="sm"
                key={item.name}
                variant={"outline-primary"}
                value={item.name}
                id={item.name}
                active={item.name === type?.name}
                onClick={() => {
                    setType(item as any)
                }}
            >
                {item.name}
            </ToggleButton>,
        )
    })

    return (
        <ToggleButtonGroup
            name={"type"}
            type={"radio"}
            style={{ paddingRight: "5px" }}
        >
            {buttons}
        </ToggleButtonGroup>
    )
}

function getCameraButtons(
    cameras: Camera[] | null,
    camera: Camera | null,
    setCamera: Dispatch<SetStateAction<null>>,
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
                variant={"outline-primary"}
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
            style={{ paddingRight: "5px" }}
        >
            {buttons}
        </ToggleButtonGroup>
    )
}

function getStreamButton(camera: Camera | null) {
    return (
        <ButtonGroup style={{ paddingRight: "5px" }}>
            <Button
                size="sm"
                variant={"outline-primary"}
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
) {
    if (!dates?.length) {
        return null
    }

    const eventDates: JSX.Element[] = []

    dates.forEach((item) => {
        const dateFriendly = item.format("YYYY-MM-DD")
        eventDates.push(
            <Dropdown.Item
                variant={"outline-primary"}
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
        <Dropdown style={{ paddingRight: "5px", width: "100%" }}>
            <Dropdown.Toggle
                variant={"outline-primary"}
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

function getObjectFilter(setObjectFilter: Dispatch<SetStateAction<string>>) {
    return (
        <FormControl
            id="objectFilter"
            type="text"
            size={"sm"}
            style={{ width: "200px" }}
            onChange={(event) => {
                setObjectFilter(event.target.value)
            }}
        />
    )
}

export interface MenuProps {
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
}

export function Menu(props: MenuProps) {
    return (
        <Navbar bg="light" expand="lg" style={{ fontSize: "10pt" }}>
            <Navbar.Brand
                href="#"
                style={{
                    fontSize: "14pt",
                    marginLeft: "10px",
                    marginRight: "10px",
                }}
                onClick={() => {}}
            >
                Cameranator
            </Navbar.Brand>

            <Navbar.Toggle />

            <Navbar.Collapse>
                <Nav>
                    {/* {getTypeButtons(props.types, props.type, props.setType)} */}
                    {getCameraButtons(
                        props.cameras,
                        props.camera,
                        props.setCamera,
                    )}
                    {/* {getStreamButton(props.camera)} */}
                    {getDateDropdown(props.dates, props.date, props.setDate)}
                </Nav>
                {getObjectFilter(props.setObjectFilter)}
            </Navbar.Collapse>
        </Navbar>
    )
}
