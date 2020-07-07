package main

import (
	"time"
)

type Profile struct {
	ProfileID   ID     `json:"id"`
	ExternalID  ID     `json:"external_id,omitempty"`
	Type       string         `json:"type"`
	Service    string         `json:"service"`
	Labels     Labels `json:"labels,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
}

func ProfileFromProfileMeta(meta Meta) Profile {
	return Profile{
		ProfileID:  meta.ProfileID,
		ExternalID: meta.ExternalID,
		Type:       meta.Type.String(),
		Service:    meta.Service,
		Labels:     meta.Labels,
		CreatedAt:  meta.CreatedAt.Truncate(time.Second),
	}
}

