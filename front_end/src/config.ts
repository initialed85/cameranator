let protocol = window?.location?.protocol || "http:";
let host = window?.location?.host || "cameranator.initialed85.cc";

if (window?.location?.port === "3000") {
  host = "cameranator.chronos";
}

const apiUrlPath = `/api/v1/graphql/`;

export const apiHttpUrl = `${protocol}//${host}${apiUrlPath}`;

export const apiWsUrl = `${
  protocol === "https:" ? "wss" : "ws"
}://${host}${apiUrlPath}`;

export const fileHttpUrl = `${protocol}//${host}/`;
