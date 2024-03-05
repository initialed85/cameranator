export interface Type {
    name: string
    is_segment: boolean | null
    is_stream: boolean
}

export const TYPES: Type[] = [
    // TODO: disabled for now
    // {
    //   name: "Event",
    //   is_segment: false,
    //   is_stream: false,
    // },
    {
        name: "Segment",
        is_segment: true,
        is_stream: false,
    },
]
