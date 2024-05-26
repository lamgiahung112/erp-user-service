package utils

import (
	"log"

	"github.com/ip2location/ip2location-go/v9"
)

type IpLocationUtils struct{}

type IpLocationData struct {
	IP      string `json:"IP"`
	Country string `json:"Country"`
	Region  string `json:"Region"`
	City    string `json:"City"`
}

var db *ip2location.DB

func InitIpLocationUtils() {
	_db, err := ip2location.OpenDB("/app/IP2LOCATION-LITE-DB3.BIN")

	if err != nil {
		log.Println(err)
		return
	}

	db = _db
}

func (*IpLocationUtils) GetLocationDatafromIP(ip string) (*IpLocationData, error) {
	result, err := db.Get_all(ip)

	if err != nil {
		return nil, err
	}

	return &IpLocationData{
		IP:      ip,
		Country: result.Country_short,
		Region:  result.Region,
		City:    result.City,
	}, nil
}
