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
    getTypeDropdown() {
        return (
            <TypeDropdown
                changeHandler={(type) => {
                    this.props.typeChangeHandler(type);
                }}
            />
        );
    }

    getDateDropdown() {
        return (
            <DateDropdown
                dates={this.props.dates}
                changeHandler={(date) => {
                    this.props.dateChangeHandler(date);
                }}
            />
        );
    }

    getCameraDropdown() {
        return (
            <CameraDropdown
                cameras={this.props.cameras}
                changeHandler={(camera) => {
                    this.props.cameraChangeHandler(camera);
                }}
            />
        );
    }

    getNavbar() {
        return (
            <Navbar bg="light" expand="lg" style={{ fontSize: "10pt" }}>
                <Navbar.Brand
                    href="#home"
                    style={{ fontSize: "14pt" }}
                    onClick={() => {
                        this.props.typeChangeHandler(null);
                        this.props.dateChangeHandler(null);
                        this.props.cameraChangeHandler(null);
                    }}
                >
                    Cameranator
                </Navbar.Brand>

                <Navbar.Toggle aria-controls="basic-navbar-nav" />

                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="mr-auto">
                        {this.getTypeDropdown()}
                        {this.getDateDropdown()}
                        {this.getCameraDropdown()}
                    </Nav>
                </Navbar.Collapse>
            </Navbar>
        );
    }

    getAlert() {
        return (
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
                variant={this.props.connected ? "success" : "danger"}
            >
                {this.props.connected ? "Online" : "Offline"}
            </Alert>
        );
    }

    getBreadcrumb() {
        return (
            <Breadcrumb style={{ fontSize: "10pt" }}>
                <Breadcrumb.Item active>{this.getAlert()}</Breadcrumb.Item>
                <Breadcrumb.Item active>
                    {getFriendlyStringForType(this.props.type)}
                </Breadcrumb.Item>
                <Breadcrumb.Item active>
                    {this.props.date
                        ? this.props.date.local().format("YYYY-MM-DD")
                        : "No date"}
                </Breadcrumb.Item>
                <Breadcrumb.Item active>
                    {this.props.camera ? this.props.camera.name : "No camera"}
                </Breadcrumb.Item>
            </Breadcrumb>
        );
    }

    getEventTable() {
        return <EventTable events={this.props.events} />;
    }

    render() {
        return (
            <Container style={{ width: "100%" }}>
                {this.getNavbar()}
                {this.getBreadcrumb()}
                {this.getEventTable()}
            </Container>
        );
    }
}

export default App;
