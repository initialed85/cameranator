export interface Image {
    uuid: string;
    file_path: string;
}

export function getImage(item: any): Image {
    return {
        uuid: item["uuid"],
        file_path: item["file_path"],
    };
}
