import {getClient, handleResultPromise} from "../utils";
import { gql } from "@apollo/client";

export interface Camera {
    uuid: string;
    name: string;
}

export interface GetCamerasHandler {
    (cameras: Camera[] | null): void;
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

    handleResultPromise(
        "camera",
        client
            .query({
                query: query,
            }),
        (data: any | null) => {
            if (data === null) {
                handler(null);
                return;
            }

            let cameras: Camera[] = [];

            data.forEach((item: any) => {
                cameras.push(getCamera(item));
            });

            handler(cameras);
        }
    )
}
