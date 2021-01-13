let possibleUri = "/api/v1/graphql";
let possibleUrlPrefix = "/";

// TODO: fix hard-coding for my home setup
if (!process.env.NODE_ENV || process.env.NODE_ENV === "development") {
    possibleUri = "http://192.168.137.253:81/api/v1/graphql";
    possibleUrlPrefix = "http://192.168.137.253:81/";
}

export const uri = possibleUri;
export const urlPrefix = possibleUrlPrefix;
