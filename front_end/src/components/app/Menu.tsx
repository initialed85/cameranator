import "bootstrap/dist/css/bootstrap.min.css";
import React, { Dispatch, SetStateAction } from "react";
import "./App.css";
import { Camera } from "../../hasura/camera";
import {
  Button,
  ButtonGroup,
  Dropdown,
  Nav,
  Navbar,
  ToggleButton,
  ToggleButtonGroup,
} from "react-bootstrap";
import moment from "moment";
import { Type } from "../../hasura/type";
import { fileHttpUrl } from "../../config";

function getTypeButtons(
  types: Type[] | null,
  type: Type | null,
  setType: Dispatch<SetStateAction<null>>
) {
  if (!types?.length) {
    return null;
  }

  const buttons: JSX.Element[] = [];

  types.map((item) => {
    buttons.push(
      <ToggleButton
        size="sm"
        key={item.name}
        variant={"outline-primary"}
        value={item.name}
        id={item.name}
        active={item.name === type?.name}
        onClick={() => {
          setType(item as any);
        }}
      >
        {item.name}
      </ToggleButton>
    );
  });

  return (
    <ToggleButtonGroup
      name={"type"}
      type={"radio"}
      style={{ paddingRight: "5px" }}
    >
      {buttons}
    </ToggleButtonGroup>
  );
}

function getStreamButton(camera: Camera | null) {
  return (
    <ButtonGroup style={{ paddingRight: "5px" }}>
      <Button
        size="sm"
        variant={"outline-primary"}
        disabled={!camera}
        href={
          camera
            ? `${fileHttpUrl}motion-stream/${camera.external_id}/stream/`
            : "#"
        }
        target={"_stream"}
      >
        Stream
      </Button>
    </ButtonGroup>
  );
}

function getCameraButtons(
  cameras: Camera[] | null,
  camera: Camera | null,
  setCamera: Dispatch<SetStateAction<null>>
) {
  if (!cameras?.length) {
    return null;
  }

  const buttons: JSX.Element[] = [];

  cameras.map((item) => {
    buttons.push(
      <ToggleButton
        size="sm"
        key={item.uuid}
        variant={"outline-primary"}
        value={item.uuid}
        id={item.uuid}
        active={item.uuid === camera?.uuid}
        onClick={() => {
          setCamera(item as any);
        }}
      >
        {item.name}
      </ToggleButton>
    );
  });

  return (
    <ToggleButtonGroup
      name={"camera"}
      type={"radio"}
      style={{ paddingRight: "5px" }}
    >
      {buttons}
    </ToggleButtonGroup>
  );
}

function getDateDropdown(
  dates: moment.Moment[] | null,
  date: moment.Moment | null,
  setDate: Dispatch<SetStateAction<null>>
) {
  if (!dates?.length) {
    return null;
  }

  const eventDates: JSX.Element[] = [];

  dates.map((item) => {
    const dateFriendly = item.format("YYYY-MM-DD");
    eventDates.push(
      <Dropdown.Item
        variant={"outline-primary"}
        href={"#"}
        key={dateFriendly}
        id={dateFriendly}
        eventKey={dateFriendly}
        active={item.toISOString() === date?.toISOString()}
        onClick={() => {
          setDate(item as any);
        }}
      >
        {dateFriendly}
      </Dropdown.Item>
    );
  });

  return (
    <Dropdown
      style={{ paddingRight: "5px" }}
      as={ButtonGroup}
      variant={"outline-primary"}
    >
      <Dropdown.Toggle
        variant={"outline-primary"}
        id="date"
        size="sm"
        // active={!!date}
      >
        {date ? date.format("YYYY-MM-DD") : "Date"}
      </Dropdown.Toggle>

      <Dropdown.Menu variant={"outline-primary"}>{eventDates}</Dropdown.Menu>
    </Dropdown>
  );
}

export interface MenuProps {
  cameras: Camera[] | null;
  camera: Camera | null;
  setCamera: Dispatch<SetStateAction<null>>;
  dates: moment.Moment[] | null;
  date: moment.Moment | null;
  setDate: Dispatch<SetStateAction<null>>;
  types: Type[] | null;
  type: Type | null;
  setType: Dispatch<SetStateAction<null>>;
}

export function Menu(props: MenuProps) {
  return (
    <Navbar bg="light" expand="lg" style={{ fontSize: "10pt" }}>
      <Navbar.Brand
        href="#"
        style={{ fontSize: "14pt", marginLeft: "10px", marginRight: "10px" }}
        onClick={() => {}}
      >
        Cameranator
      </Navbar.Brand>

      <Navbar.Toggle aria-controls="basic-navbar-nav" />

      <Navbar.Collapse>
        <Nav>
          {getTypeButtons(props.types, props.type, props.setType)}
          {getStreamButton(props.camera)}
          {getCameraButtons(props.cameras, props.camera, props.setCamera)}
          {getDateDropdown(props.dates, props.date, props.setDate)}
        </Nav>
      </Navbar.Collapse>
    </Navbar>
  );
}
