const httpProtocol = window?.location?.protocol || "http:";
const wsProtocol = httpProtocol === "https:" ? "wss:" : "ws:";

const host =
  window?.location?.port === "3000"
    ? "cameranator.initialed85.cc"
    : window?.location?.host;

const apiUrlPath = `/api/v1/graphql`;

export const apiHttpUrl = `${httpProtocol}//${host}${apiUrlPath}/`; // trailing slash

export const apiWsUrl = `${wsProtocol}//${host}${apiUrlPath}`;

export const fileHttpUrl = `${httpProtocol}//${host}/`;
