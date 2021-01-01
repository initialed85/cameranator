import { ApolloClient, InMemoryCache } from "@apollo/client";

// TODO: make this configurable
const uri = "cameranator/api/v1/graphql";

export function getClient() {
    return new ApolloClient({
        uri: uri,
        cache: new InMemoryCache(),
    });
}
