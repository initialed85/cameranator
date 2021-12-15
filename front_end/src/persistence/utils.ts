import { ApolloClient, InMemoryCache } from "@apollo/client";
import { uri } from "../config/config";

let client: undefined | ApolloClient<any>;

export function getClient(): ApolloClient<any> {
    if (!client) {
        client = new ApolloClient({
            uri: uri,
            cache: new InMemoryCache(),
        });
    }

    return client;
}
