import React from "react";

import "./App.css";
import { Container, Nav, Navbar } from "react-bootstrap";
import { Camera, getCameras } from "./persistence/camera";
import { CameraDropdown } from "./components/camera_drop_down";
import { getDates } from "./persistence/date";
import moment from "moment/moment";
import { DateDropdown } from "./components/date_drop_down";
import { Event, getEvents } from "./persistence/event";
import { EventTable } from "./components/event_table";
import { SEGMENTS, TypeDropdown } from "./components/type_drop_down";

interface AppProps {}

interface AppState {
    dates: moment.Moment[];
    cameras: Camera[];
    events: Event[];
    type: string | null;
    date: moment.Moment | null;
    camera: Camera | null;
}

class App extends React.Component<AppProps, AppState> {
    constructor(props: AppProps, state: AppState) {
        super(props, state);

        this.state = {
            dates: [] as moment.Moment[],
            cameras: [] as Camera[],
            events: [] as Event[],
            type: null,
            date: null,
            camera: null,
        };

        setInterval(() => {
            this.update()
        }, 1000);
    }

    update() {
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
        getDates((dates) => {
            this.setState({ dates: dates });
        });

        getCameras((cameras) => {
            this.setState({ cameras: cameras });
        });
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

                <EventTable events={this.state.events} />
            </Container>
        );
    }
}

export default App;
