import { gql } from "@apollo/client";

export const CAMERAS = gql`
  query cameras {
    camera(order_by: { name: asc }) {
      uuid
      name
      stream_url
      external_id
    }
  }
`;

export interface Camera {
  __typename: string | null;
  uuid: string;
  name: string;
  stream_url: string | null;
  external_id: string | null;
}
