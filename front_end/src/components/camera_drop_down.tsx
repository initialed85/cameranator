import { Camera } from "../persistence/camera";
import React from "react";
import { NavDropdown } from "react-bootstrap";

export interface CameraDropDownChangeHandler {
    (camera: Camera): void;
}

interface CameraDropdownProps {
    cameras: Camera[];
    changeHandler: CameraDropDownChangeHandler;
}

export class CameraDropdown extends React.Component<CameraDropdownProps, any> {
    render() {
        let items: any[] = [];

        this.props.cameras.forEach((camera: Camera) => {
            items.push(
                <NavDropdown.Item
                    href={`#camera/${camera.uuid}`}
                    key={camera.uuid}
                    onClick={() => {
                        this.props.changeHandler(camera);
                    }}
                >
                    {camera.name}
                </NavDropdown.Item>
            );
        });

        return (
            <NavDropdown title="Camera" id="camera-dropdown">
                {items}
            </NavDropdown>
        );
    }
}
