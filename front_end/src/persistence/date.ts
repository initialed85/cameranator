import { getClient } from "./utils";
import { DocumentNode, gql } from "@apollo/client";
import moment from "moment/moment";

export interface GetDatesHandler {
    (dates: moment.Moment[]): void;
}

function getQuery(): DocumentNode {
    return gql(`
{
  event(
    order_by: {start_timestamp: desc}
  ) {
    start_timestamp
  }
}
`);
}

export function getDate(item: any): moment.Moment {
    return moment.utc(item["start_timestamp"]);
}

export function getDates(
    handler: GetDatesHandler
) {
    const client = getClient();

    client
        .query({
            query: getQuery(),
        })
        .catch((e) => {
            console.error(e);
        })
        .then((r) => {
            const data = (r as any).data["event"].slice();

            let dates: moment.Moment[] = [];
            data.forEach((item: any) => {
                dates.push(getDate(item).local());
            });

            let dateByShortDate = new Map<string, moment.Moment>();
            dates.forEach((date) => {
                dateByShortDate.set(date.format("YYYY-MM-DD"), date);
            });

            const deduplicatedDates = Array.from(dateByShortDate.values());

            handler(deduplicatedDates);
        });
}
