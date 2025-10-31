package monitor

import "time"

// ResourceHealth holds health status for a resource
type ResourceHealth struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Type              string    `json:"type"`
	AvailabilityState string    `json:"availabilityState"`
	Summary           string    `json:"summary,omitempty"`
	ReasonType        string    `json:"reasonType,omitempty"`
	LastUpdated       time.Time `json:"lastUpdated"`
	Healthy           bool      `json:"healthy"`
}
