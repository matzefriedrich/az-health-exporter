package monitor

import "time"

// ResourceHealth holds health status for a resource
type ResourceHealth struct {
	ResourceID        string    `json:"resource_id"`
	ResourceName      string    `json:"resource_name"`
	ResourceType      string    `json:"resource_type"`
	AvailabilityState string    `json:"availability_state"`
	Summary           string    `json:"summary"`
	ReasonType        string    `json:"reason_type"`
	LastUpdated       time.Time `json:"last_updated"`
	Healthy           bool      `json:"healthy"`
}
