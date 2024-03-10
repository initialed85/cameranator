import { gql } from "@apollo/client"

export const CAMERAS = gql`
    subscription cameras {
        camera(order_by: { name: asc }) {
            id
            name
            stream_url
        }
    }
`

export interface Camera {
    __typename: string | null
    id: string
    name: string
    stream_url: string | null
}
