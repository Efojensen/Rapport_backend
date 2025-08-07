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

type GeoApifyResponse struct {
	Results []struct {
		Properties struct {
			Name     string  `json:"name"`
			Country  string  `json:"country"`
			Suburb   string  `json:"suburb"`
			Street   string  `json:"street"`
			Distance float64 `json:"distance"`
			City     string  `json:"city"`
			Address  string  `json:"formatted"`
		}`json:"properties"`
	}`json:"features"`
}

func (report *SOSReport) GetLatLongAddress() (*GeoApifyResponse, error) {
	geoApi := os.Getenv("GEOAPIFY_KEY")
	lat := strconv.FormatFloat(report.GeoLocation.Latitude, 'f', -1, 64)
	long := strconv.FormatFloat(report.GeoLocation.Longitude, 'f', -1, 64)

	httpMethod := fmt.Sprintf("https://api.geoapify.com/v1/geocode/reverse?lat=%s&lon=%s&apiKey=%s", lat, long, geoApi)

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Get(httpMethod)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var SosLoc GeoApifyResponse
	err = json.Unmarshal(body, &SosLoc)

	if err != nil {
		return nil, err
	}

	if len(SosLoc.Results) == 0 {
		return nil, fmt.Errorf("no location features found")
	}

	return &SosLoc, nil
}
