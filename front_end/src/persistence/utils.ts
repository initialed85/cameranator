import { ApolloClient, InMemoryCache } from "@apollo/client";

// TODO: make this configurable
const uri = "http://192.168.137.253:8082/v1/graphql";

export function getClient() {
    return new ApolloClient({
        uri: uri,
        cache: new InMemoryCache(),
    });
}
