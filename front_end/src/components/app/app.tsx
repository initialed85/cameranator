import React from "react";

import "./app.css";
import { Alert, Breadcrumb, Container, Nav, Navbar } from "react-bootstrap";
import { CameraDropdown } from "../camera_drop_down/camera_drop_down";
import { DateDropdown } from "../date_drop_down/date_drop_down";
import { EventTable } from "../event_table/event_table";
import {
    EVENTS,
    SEGMENTS,
    TypeDropdown,
} from "../type_drop_down/type_drop_down";
import { AppProps } from "./app_props";

function getFriendlyStringForType(type: string | null): string {
    if (type === EVENTS) {
        return "Events";
    } else if (type === SEGMENTS) {
        return "Segments";
    } else {
        return "No type";
    }
}

class App extends React.Component<AppProps> {
    render() {
        return (
            <Container style={{ width: "100%" }}>
                <Navbar bg="light" expand="lg" style={{ fontSize: "10pt" }}>
                    <Navbar.Brand href="#home" style={{ fontSize: "14pt" }}>
                        Cameranator
                    </Navbar.Brand>

                    <Navbar.Toggle aria-controls="basic-navbar-nav" />

                    <Navbar.Collapse id="basic-navbar-nav">
                        <Nav className="mr-auto">
                            <TypeDropdown
                                changeHandler={(type) => {
                                    this.props.typeChangeHandler(type);
                                }}
                            />

                            <DateDropdown
                                dates={this.props.dates}
                                changeHandler={(date) => {
                                    this.props.dateChangeHandler(date);
                                }}
                            />

                            <CameraDropdown
                                cameras={this.props.cameras}
                                changeHandler={(camera) => {
                                    this.props.cameraChangeHandler(camera);
                                }}
                            />
                        </Nav>
                    </Navbar.Collapse>
                </Navbar>

                <Breadcrumb style={{ fontSize: "10pt" }}>
                    <Breadcrumb.Item active>
                        <Alert
                            style={{
                                paddingTop: 0,
                                paddingBottom: 0,
                                paddingLeft: 5,
                                paddingRight: 5,
                                margin: 0,
                                width: 75,
                                textAlign: "center",
                            }}
                            variant={
                                this.props.connected ? "success" : "danger"
                            }
                        >
                            {this.props.connected ? "Online" : "Offline"}
                        </Alert>
                    </Breadcrumb.Item>
                    <Breadcrumb.Item active>
                        {getFriendlyStringForType(this.props.type)}
                    </Breadcrumb.Item>
                    <Breadcrumb.Item active>
                        {this.props.date
                            ? this.props.date.local().format("YYYY-MM-DD")
                            : "No date"}
                    </Breadcrumb.Item>
                    <Breadcrumb.Item active>
                        {this.props.camera
                            ? this.props.camera.name
                            : "No camera"}
                    </Breadcrumb.Item>
                </Breadcrumb>

                <EventTable events={this.props.events} />
            </Container>
        );
    }
}

export default App;
