package libtools

import (
	"fmt"
	"math"
	"strings"
)

// 返回单位为：米
func GetDistance(lat1, lng1, lat2, lng2 float64) float64 {
	radius := 6371000.0 //6378137.0
	rad := math.Pi / 180.0
	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius
}

func BuildEsGeoLocation(lng, lat string) string {
	return fmt.Sprintf(`%s, %s`, lat, lng)
}

func EsGeoLocation2LngLat(location string) (lng, lat string) {
	exp := strings.Split(location, ",")
	if len(exp) == 2 {
		lat = strings.TrimSpace(exp[0])
		lng = strings.TrimSpace(exp[1])
	}

	return
}

// caculateTimeZone计算时区
func CaculateTimeZone(lon float64) string {
	var timeZone float64
	shangValue := (lon / 15)
	yushuValue := math.Abs(math.Mod(lon, 15))
	if yushuValue <= 7.5 {
		timeZone = shangValue
	} else {
		if lon > 0 {
			timeZone = shangValue + 1
		} else {
			timeZone = shangValue - 1
		}

	}
	return fmt.Sprintf("%.1f", timeZone)
}
