import { getClient } from "../utils";
import { gql } from "@apollo/client";

export interface Camera {
    uuid: string;
    name: string;
}

export interface GetCamerasHandler {
    (cameras: Camera[]): void;
}

const query = gql(`
query {
  camera(order_by: {name: asc}) {
    uuid
    name
  }
}
`);

export function getCamera(item: any): Camera {
    return {
        uuid: item["uuid"],
        name: item["name"],
    };
}

export function getCameras(handler: GetCamerasHandler) {
    const client = getClient();

    client
        .query({
            query: query,
        })
        .catch((e) => {
            console.error(e);
        })
        .then((r) => {
            let cameras: Camera[] = [];

            const data = (r as any).data["camera"].slice();

            data.forEach((item: any) => {
                cameras.push(getCamera(item));
            });

            handler(cameras);
        });
}
