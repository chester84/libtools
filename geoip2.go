// https://dev.maxmind.com/geoip/geoip2/geolite2/
// https://github.com/oschwald/geoip2-golang

package libtools

import (
	"log"
	"net"

	"github.com/beego/beego/v2/core/config"
	"github.com/beego/beego/v2/core/logs"
	"github.com/oschwald/geoip2-golang"
)

type Location struct {
	Latitude  float64
	Longitude float64
	TimeZone  string
}

const EmptyRecord string = "-"

var geoIpDB *geoip2.Reader

func init() {
	geoip2DbName, _ := config.String("geolite2_city_dbname")
	db, err := geoip2.Open(geoip2DbName)
	if err != nil {
		logs.Error("wrong config geolite2_city_dbname: ", geoip2DbName, ", err:", err)
		log.Fatal("can not init geoip2")
	}

	geoIpDB = db
}

func getGeoipCityDb(ipOrigin string) (*geoip2.City, error) {
	ip := net.ParseIP(ipOrigin)
	record, err := geoIpDB.City(ip)
	if err != nil {
		logs.Error("Can not find record. ip:", ipOrigin, ", err: ", err)
		return nil, err
	}

	return record, err
}

// ip取ISO国家码
func GeoipISOCountryCode(ipOrigin string) string {
	record, err := getGeoipCityDb(ipOrigin)
	if err != nil {
		return EmptyRecord
	}

	if record.Country.IsoCode == "" {
		return EmptyRecord
	}

	return record.Country.IsoCode
}

func getCity(ipOrigin, lang string) string {
	record, err := getGeoipCityDb(ipOrigin)
	if err != nil {
		return EmptyRecord
	}

	if record.City.Names[lang] == "" {
		return EmptyRecord
	}

	return record.City.Names[lang]
}

func GeoIpCityEn(ipOrigin string) string {
	return getCity(ipOrigin, "en")
}

func GeoIpCityZhCN(ipOrigin string) string {
	return getCity(ipOrigin, "zh-CN")
}

func GeoIpLocation(ipOrigin string) Location {
	record, err := getGeoipCityDb(ipOrigin)
	var l Location
	if err == nil {
		l.Latitude = record.Location.Latitude
		l.Longitude = record.Location.Longitude
		l.TimeZone = record.Location.TimeZone
	}

	return l
}
