import Collection from "../collection";
import {DocumentNode, gql} from "@apollo/client";
import moment from "moment/moment";
import {info} from "../../common/utils";

function getQuery(args: any): DocumentNode {
    return gql(`
{
  event(
    order_by: {start_timestamp: desc}
    distinct_on: start_timestamp
  ) {
    start_timestamp
  }
}
`);
}

export function getDate(item: any): moment.Moment {
    return moment.utc(item["start_timestamp"]);
}

export class DateCollection extends Collection {
    constructor() {
        super(getQuery, "event");
    }

    get(args: any): Promise<any> {
        info(`${this.constructor.name}.get fired`);
        return new Promise((resolve, reject) => {
            this.handleResultPromise(this.getResultPromise(args))
                .catch((e) => {
                    reject(e);
                })
                .then((data) => {
                    let dates: moment.Moment[] = [];
                    data.forEach((item: any) => {
                        dates.push(getDate(item).local());
                    });

                    let dateByShortDate = new Map<string, moment.Moment>();
                    dates.forEach((date) => {
                        dateByShortDate.set(date.format("YYYY-MM-DD"), date);
                    });

                    const deduplicatedDates = Array.from(
                        dateByShortDate.values()
                    );

                    resolve(deduplicatedDates);
                });
        });
    }
}

export default DateCollection;
