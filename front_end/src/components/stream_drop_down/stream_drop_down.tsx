import { Camera } from "../../persistence/collections/camera";
import React from "react";
import { NavDropdown } from "react-bootstrap";
import { urlPrefix } from "../../config/config";

function getMJPEGPath(camera: Camera): string {
    return `${urlPrefix}motion-stream/${camera.external_id}/stream`;
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
                    rel={"noreferrer"}
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
