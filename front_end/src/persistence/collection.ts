import { ApolloClient, ApolloQueryResult } from "@apollo/client";
import { getClient } from "./utils";

export abstract class Collection {
    client: ApolloClient<any>;
    key: string;
    getQuery: any;

    protected constructor(getQuery: CallableFunction, key: string) {
        this.getQuery = getQuery;
        this.key = key;
        this.client = getClient();
    }

    getResultPromise(args: any): Promise<ApolloQueryResult<any>> {
        return this.client.query({
            query: this.getQuery(args),
        });
    }

    handleResultPromise(
        resultPromise: Promise<ApolloQueryResult<any>>
    ): Promise<any> {
        return new Promise((resolve, reject) => {
            resultPromise
                .catch((e) => {
                    console.warn(
                        `warning: attempt to get ${this.key} caused: `,
                        e
                    );
                    reject(e);
                })
                .then((r) => {
                    if (!r) {
                        reject(new Error("no result"));
                        return;
                    }

                    if (!r.data) {
                        reject(new Error("no result.data"));
                        return;
                    }

                    const data = (r as any).data[this.key].slice();
                    if (!data) {
                        reject(`no result.data[${this.key}]`);
                        return;
                    }

                    resolve(data);
                });
        });
    }

    abstract get(args: any): Promise<any>;
}

export default Collection;
