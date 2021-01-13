import { Camera } from "../../persistence/types/camera";
import React from "react";
import { NavDropdown } from "react-bootstrap";

function getMJPEGPath(camera: Camera): string {
    return `/motion-stream/${camera.external_id}/stream`;
}

interface StreamDropdownProps {
    cameras: Camera[];
}

export class StreamDropdown extends React.Component<StreamDropdownProps, any> {
    render() {
        let items: any[] = [];

        this.props.cameras.forEach((camera: Camera) => {
            items.push(
                <NavDropdown.Item
                    style={{ fontSize: "10pt" }}
                    href={getMJPEGPath(camera)}
                    target={"_blank"}
                    key={camera.uuid}
                >
                    {camera.name}
                </NavDropdown.Item>
            );
        });

        return (
            <NavDropdown title="Stream" id="camera-dropdown">
                {items}
            </NavDropdown>
        );
    }
}
