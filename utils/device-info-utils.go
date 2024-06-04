package utils

import (
	"errors"
	"time"

	"github.com/mssola/useragent"
)

type DeviceInfo struct {
	Model      string    `json:"model"`
	OS         string    `json:"os"`
	Browser    string    `json:"browser"`
	IP         string    `json:"IP"`
	Country    string    `json:"country"`
	Region     string    `json:"region"`
	City       string    `json:"city"`
	LoggedInAt time.Time `json:"loggedInAt"`
}

type DeviceInfoUtilities struct {
}

func (d1 *DeviceInfo) Compare(d2 *DeviceInfo) bool {
	isSameBrowser := d1.Browser == d2.Browser
	isSameModel := d1.Model == d2.Model
	isSameOs := d1.OS == d2.OS

	return isSameBrowser && isSameModel && isSameOs
}

func (u *DeviceInfoUtilities) GetDevice(userAgent string, iplocation *IpLocationData) (*DeviceInfo, error) {
	result := useragent.New(userAgent)
	browser, _ := result.Browser()

	if result.Bot() {
		return nil, errors.New("bot detected")
	}

	return &DeviceInfo{
		Model:   result.Model(),
		OS:      result.OSInfo().Name,
		Browser: browser,
		IP:      iplocation.IP,
		Country: iplocation.Country,
		Region:  iplocation.Region,
		City:    iplocation.City,
	}, nil
}
