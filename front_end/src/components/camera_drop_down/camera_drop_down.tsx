import {Camera} from "../../persistence/collections/camera";
import React from "react";
import {NavDropdown} from "react-bootstrap";

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
                    style={{fontSize: "10pt"}}
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
