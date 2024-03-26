WITH
    detections_1 AS (
        SELECT
            d.*
        FROM
            detections d
        WHERE
            d.timestamp >= '2024-03-18T18:17:00+08:00'
            AND d.timestamp <= '2024-03-18T18:18:00+08:00'
            -- AND d.class_name IN ('person', 'cat', 'car', 'motorcycle')
            AND d.camera_id = 1
            AND d.class_name IN ('person')
            AND d.score > 0.55
    ),
    detections_2 AS (
        SELECT
            d1.camera_id,
            d1.timestamp,
            d1.class_name
        FROM
            detections_1 d1
    ),
    detections_3 AS (
        SELECT
            d1.*,
            st_area (d1.bounding_box::geometry) / (1920.0 * 1080.0) AS area,
            (st_xmax (d1.bounding_box::geometry) - st_xmin (d1.bounding_box::geometry)) / (1920.0 * 1080.0) AS width,
            (st_ymax (d1.bounding_box::geometry) - st_ymin (d1.bounding_box::geometry)) / (1920.0 * 1080.0) AS height,
            (st_xmax (d1.bounding_box::geometry) - st_xmin (d1.bounding_box::geometry)) / (st_ymax (d1.bounding_box::geometry) - st_ymin (d1.bounding_box::geometry)) AS aspect_ratio
        FROM
            detections_2 d2
            INNER JOIN detections_1 d1 ON d1.camera_id = d2.camera_id
            AND d1.timestamp = d2.timestamp
            AND d1.class_name = d2.class_name
        ORDER BY
            d1.timestamp ASC
    ),
    detections_4 AS (
        SELECT
            d3a.id,
            d3a.camera_id,
            d3a.timestamp,
            d3a.class_name,
            round(d3a.area::numeric, 3),
            d3b.id AS other_id,
            d3b.camera_id AS other_camera_id,
            d3b.timestamp AS other_timestamp,
            d3b.class_name AS other_class_name,
            round(d3b.area::numeric, 3) AS other_area,
            round(
                (
                    (
                        extract(
                            epoch
                            from
                                d3b.timestamp
                        ) - extract(
                            epoch
                            from
                                d3a.timestamp
                        )
                    ) / 60.0
                )::numeric,
                3
            ) AS time_factor,
            round((d3b.area / d3a.area)::numeric, 3) AS area_factor,
            round((d3b.aspect_ratio / d3a.aspect_ratio)::numeric, 3) AS aspect_ratio_factor,
            round((d3b.width / d3a.width)::numeric, 3) AS width_factor,
            round((d3b.height / d3a.height)::numeric, 3) AS height_factor,
            round((st_distance (d3b.centroid::geometry, d3a.centroid::geometry) / (1920.0 * 1080.0))::numeric, 3) AS distance_factor
        FROM
            detections_3 d3a
            INNER JOIN LATERAL (
                SELECT
                    d3b.*
                FROM
                    detections_3 d3b
                WHERE
                    d3b.camera_id = d3a.camera_id
                    AND d3b.class_name = d3a.class_name
                    AND d3b.timestamp > d3a.timestamp
            ) AS d3b ON true
    ),
    detections_5 AS (
        SELECT
            d4.*
        FROM
            detections_4 d4
        WHERE
            time_factor <= 0.005
            AND area_factor >= 0.99
            AND area_factor <= 1.01
            AND aspect_ratio_factor >= 0.99
            AND aspect_ratio_factor <= 1.01
            AND distance_factor <= 0.001
    )
SELECT
    *
FROM
    detections_5;
