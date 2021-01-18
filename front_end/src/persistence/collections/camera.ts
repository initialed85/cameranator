import { gql } from "@apollo/client";
import Collection from "../collection";

function getQuery(args: any) {
    return gql(`
query {
  camera(order_by: {name: asc}) {
    uuid
    name
    stream_url
    external_id
  }
}
`);
}

export interface Camera {
    uuid: string;
    name: string;
    stream_url: string;
    external_id: string;
}

export function getCamera(item: any): Camera {
    return {
        uuid: item["uuid"],
        name: item["name"],
        stream_url: item["stream_url"],
        external_id: item["external_id"],
    };
}

class CameraCollection extends Collection {
    constructor() {
        super(getQuery, "camera");
    }

    get(args: any): Promise<any> {
        return new Promise((resolve, reject) => {
            this.handleResultPromise(this.getResultPromise(args))
                .catch((e) => {
                    reject(e);
                })
                .then((data) => {
                    let cameras: Camera[] = [];

                    data.forEach((item: any) => {
                        cameras.push(getCamera(item));
                    });

                    resolve(cameras);
                });
        });
    }
}

export default CameraCollection;
