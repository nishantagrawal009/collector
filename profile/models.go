package profile

import (
	"time"
)

type Profile struct {
	ProfileID  ID        `json:"id"`
	ExternalID ID        `json:"external_id,omitempty"`
	Type       string    `json:"type"`
	Service    string    `json:"service"`
	Labels     Labels    `json:"labels,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
}


