export interface Type {
    name: string
    is_segment: boolean | null
    is_stream: boolean
}

export const TYPES: Type[] = [
    {
        name: "Segment",
        is_segment: true,
        is_stream: false,
    },
]
