package models

import "time"

type SOSReport struct {
	Location Location  `json:"location"`
	SentTime time.Time `json:"sentTime"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Accuracy  float64 `json:"accuracy"`
}