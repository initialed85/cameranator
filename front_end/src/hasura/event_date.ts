import { gql } from "@apollo/client";
import moment from "moment/moment";

export const EVENT_DATES = gql`
  subscription dates {
    event(order_by: { start_timestamp: desc }, distinct_on: start_timestamp) {
      start_timestamp
    }
  }
`;

export interface EventDate {
  __typename: string;
  start_timestamp: string;
}

export const getDeduplicatedDates = (data: any) => {
  let dates: moment.Moment[] = [];
  (data?.event || []).forEach((item: EventDate) => {
    dates.push(moment.utc(item.start_timestamp).local());
  });

  let dateByShortDate = new Map<string, moment.Moment>();
  dates.forEach((date) => {
    dateByShortDate.set(date.format("YYYY-MM-DD"), date);
  });

  return Array.from(dateByShortDate.values());
};
