import React from "react";
import { NavDropdown } from "react-bootstrap";

export const EVENTS = "events";
export const SEGMENTS = "segments";

export interface TypeDropDownChangeHandler {
    (type: string): void;
}

interface TypeDropdownProps {
    changeHandler: TypeDropDownChangeHandler;
}

export class TypeDropdown extends React.Component<TypeDropdownProps, any> {
    render() {
        return (
            <NavDropdown title="Type" id="camera-dropdown">
                <NavDropdown.Item
                    style={{ fontSize: "10pt" }}
                    href={`#type/${EVENTS}`}
                    key={EVENTS}
                    onClick={() => {
                        this.props.changeHandler(EVENTS);
                    }}
                >
                    Events
                </NavDropdown.Item>
                <NavDropdown.Item
                    style={{ fontSize: "10pt" }}
                    href={`#type/${SEGMENTS}`}
                    key={SEGMENTS}
                    onClick={() => {
                        this.props.changeHandler(SEGMENTS);
                    }}
                >
                    Segments
                </NavDropdown.Item>
            </NavDropdown>
        );
    }
}
