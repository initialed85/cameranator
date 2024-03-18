package object_tracker

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseGeometry(rawGeometry string) ([]Point, error) {
	points := make([]Point, 0)

	for _, rawXYAsStr := range strings.Split(rawGeometry, "),(") {
		rawXYAsStr = strings.ReplaceAll(rawXYAsStr, "(", "")
		rawXYAsStr = strings.ReplaceAll(rawXYAsStr, ")", "")
		rawXYAsSlice := strings.Split(strings.Trim(rawXYAsStr, " ()"), ",")

		x, err := strconv.ParseFloat(rawXYAsSlice[0], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %#+v[0] as float", rawXYAsSlice)
		}

		y, err := strconv.ParseFloat(rawXYAsSlice[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %#+v[1] as float", rawXYAsSlice)
		}

		points = append(points, Point{X: x, Y: y})
	}

	return points, nil
}
