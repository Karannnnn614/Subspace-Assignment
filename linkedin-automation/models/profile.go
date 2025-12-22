package models

import "time"

// Profile represents a LinkedIn profile
type Profile struct {
	Name             string
	Headline         string
	ProfileURL       string
	Company          string
	Location         string
	ConnectionDegree string
	SearchKeyword    string
	DiscoveredAt     time.Time
}
