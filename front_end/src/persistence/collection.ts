import { ApolloClient, ApolloQueryResult } from "@apollo/client";
import { getClient } from "./utils";
import { info, warn } from "../common/utils";

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
        info(`${this.constructor.name}.getResultPromise fired`);
        return this.client.query({
            query: this.getQuery(args),
        });
    }

    handleResultPromise(
        resultPromise: Promise<ApolloQueryResult<any>>
    ): Promise<any> {
        info(`${this.constructor.name}.handleResultPromise fired`);
        return new Promise((resolve, reject) => {
            resultPromise
                .catch((e) => {
                    warn(
                        `${this.constructor.name}.handleResultPromise attempt to get ${this.key} caused: ${e}`
                    );
                    reject(e);
                })
                .then((r) => {
                    if (!r) {
                        warn(
                            `${this.constructor.name}.handleResultPromise had no result`
                        );
                        reject(new Error("no result"));
                        return;
                    }

                    if (!r.data) {
                        warn(
                            `${this.constructor.name}.handleResultPromise had no result.data`
                        );
                        reject(new Error("no result.data"));
                        return;
                    }

                    const data = (r as any).data[this.key].slice();
                    if (!data) {
                        warn(
                            `${this.constructor.name}.handleResultPromise had no result.data[${this.key}]`
                        );
                        reject(`no result.data[${this.key}]`);
                        return;
                    }

                    info(
                        `${this.constructor.name}.handleResultPromise resolving`
                    );
                    resolve(data);
                    info(
                        `${this.constructor.name}.handleResultPromise resolved`
                    );
                });
        });
    }

    abstract get(args: any): Promise<any>;
}

export default Collection;
