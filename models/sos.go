package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type SOSReport struct {
	GeoLocation GeoLocation `json:"location"`
	SentTime    time.Time   `json:"sentTime"`
}

type GeoLocation struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type Location struct {
	Name     string `json:"name"`
	Country  string `json:"country"`
	Suburb   string `json:"suburb"`
	Street   string `json:"street"`
	Distance int    `json:"distance"`
	City     string `json:"city"`
	Address  string `json:"formatted"`
}

func (report *SOSReport) GetLatLongAddress() (*Location, error) {
	geoApi := os.Getenv("GEOAPIFY_KEY")
	lat := strconv.FormatFloat(report.GeoLocation.Latitude, 'f', -1, 64)
	long := strconv.FormatFloat(report.GeoLocation.Longitude, 'f', -1, 64)

	httpMethod := fmt.Sprintf("https://api.geoapify.com/v1/geocode/reverse?lat=%s&lon=%s&apiKey=%s", lat, long, geoApi)

	res, err := http.Get(httpMethod)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var SosLoc Location
	err = json.Unmarshal(body, &SosLoc)

	if err != nil {
		return nil, err
	}

	return &SosLoc, nil
}
