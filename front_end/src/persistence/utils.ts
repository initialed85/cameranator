import {ApolloClient, InMemoryCache} from "@apollo/client";
import {uri} from "../config";

export function getClient() {
    return new ApolloClient({
        uri: uri,
        cache: new InMemoryCache(),
    });
}
