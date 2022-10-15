import {ApolloClient, InMemoryCache} from "@apollo/client";
import {apiUrl} from "../config/config";

let client: undefined | ApolloClient<any>;

export function getClient(): ApolloClient<any> {
    if (!client) {
        client = new ApolloClient({
            uri: apiUrl,
            cache: new InMemoryCache(),
        });
    }

    return client;
}
