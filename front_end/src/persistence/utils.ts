import { ApolloClient, ApolloQueryResult, InMemoryCache } from "@apollo/client";
import { uri } from "../config/config";

export function getClient() {
    return new ApolloClient({
        uri: uri,
        cache: new InMemoryCache(),
    });
}

export function getResultPromise(query: any): Promise<ApolloQueryResult<any>> {
    return getClient().query({
        query: query,
    });
}

export function handleResultPromise(
    key: string,
    result: Promise<ApolloQueryResult<any>>,
    handler: CallableFunction
) {
    result
        .catch((e) => {
            console.warn(`warning: attempt to get ${key} caused: `, e);
            handler(null);
        })
        .then((r) => {
            if (!r) {
                handler(null);
                return;
            }

            if (!r.data) {
                handler(null);
                return;
            }

            if (!r.data[key]) {
                handler(null);
                return;
            }

            const data = (r as any).data[key].slice();
            if (!data) {
                handler(null);
                return;
            }

            handler(data);
        });
}
