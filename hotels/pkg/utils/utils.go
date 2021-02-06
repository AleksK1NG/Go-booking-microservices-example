package utils

import (
	"fmt"
	"strings"
)

func ParsePoint(point string) (string, string) {
	longlat := strings.Split(point, ",")

	return longlat[0], longlat[1]
}

func GeneratePointToGeoFromFloat64(latitude float64, longitude float64) string {
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", latitude, longitude)
}
