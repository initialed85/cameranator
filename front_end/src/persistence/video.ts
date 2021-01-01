export interface Video {
    uuid: string;
    file_path: string;
}

export function getVideo(item: any): Video {
    return {
        uuid: item["uuid"],
        file_path: item["file_path"],
    };
}
