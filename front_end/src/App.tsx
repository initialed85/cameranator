import React from "react";

import "./App.css";
import { Alert, Breadcrumb, Container, Nav, Navbar } from "react-bootstrap";
import { Camera, getCameras } from "./persistence/camera";
import { CameraDropdown } from "./components/camera_drop_down";
import { getDates } from "./persistence/date";
import moment from "moment/moment";
import { DateDropdown } from "./components/date_drop_down";
import { Event, getEvents } from "./persistence/event";
import { EventTable } from "./components/event_table";
import { EVENTS, SEGMENTS, TypeDropdown } from "./components/type_drop_down";

interface AppProps {}

interface AppState {
    connected: boolean;
    dates: moment.Moment[];
    cameras: Camera[];
    events: Event[];
    type: string | null;
    date: moment.Moment | null;
    camera: Camera | null;
}

function getFriendlyStringForType(type: string | null): string {
    if (type === EVENTS) {
        return "Events";
    } else if (type === SEGMENTS) {
        return "Segments";
    } else {
        return "No type";
    }
}

class App extends React.Component<AppProps, AppState> {
    mounted: boolean;

    constructor(props: AppProps, state: AppState) {
        super(props, state);

        this.mounted = false;

        this.state = {
            connected: false,
            dates: [] as moment.Moment[],
            cameras: [] as Camera[],
            events: [] as Event[],
            type: null,
            date: null,
            camera: null,
        };

        setInterval(() => {
            this.update();
        }, 1000);
    }

    update() {
        if (this.mounted) {
            this.setState({ connected: !!this.state.cameras });
        }

        if (!(this.state.type && this.state.date && this.state.camera)) {
            return;
        }

        getEvents(
            this.state.type === SEGMENTS,
            this.state.date,
            this.state.camera.name,
            (events) => {
                this.setState({ events: events });
            }
        );
    }

    componentDidMount() {
        this.mounted = true;

        getDates((dates) => {
            this.setState({ dates: dates });
        });

        getCameras((cameras) => {
            this.setState({ cameras: cameras });
        });
    }

    componentWillUnmount() {
        this.mounted = false;
    }

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
                                    this.setState({ type: type }, this.update);
                                }}
                            />

                            <DateDropdown
                                dates={this.state.dates}
                                changeHandler={(date) => {
                                    this.setState({ date: date }, this.update);
                                }}
                            />

                            <CameraDropdown
                                cameras={this.state.cameras}
                                changeHandler={(camera) => {
                                    this.setState(
                                        { camera: camera },
                                        this.update
                                    );
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
                                textAlign: "center"
                            }}
                            variant={
                                this.state.connected ? "success" : "danger"
                            }
                        >
                            {this.state.connected
                                ? "Online"
                                : "Offline"}
                        </Alert>
                    </Breadcrumb.Item>
                    <Breadcrumb.Item active>
                        {getFriendlyStringForType(this.state.type)}
                    </Breadcrumb.Item>
                    <Breadcrumb.Item active>
                        {this.state.date
                            ? this.state.date.local().format("YYYY-MM-DD")
                            : "No date"}
                    </Breadcrumb.Item>
                    <Breadcrumb.Item active>
                        {this.state.camera
                            ? this.state.camera.name
                            : "No camera"}
                    </Breadcrumb.Item>
                </Breadcrumb>

                <EventTable events={this.state.events} />
            </Container>
        );
    }
}

export default App;
