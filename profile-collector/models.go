package profefe

import (
	"collector/profile"
	"time"

)

// Profile is the JSON representation of a profile returned with API response.
type Profile struct {
	ProfileID  profile.ID     `json:"id"`
	ExternalID profile.ID     `json:"external_id,omitempty"`
	Type       string         `json:"type"`
	Service    string         `json:"service"`
	Labels     profile.Labels `json:"labels,omitempty"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
}

type GetProfileDisplay struct{
	ServiceName string
	Cpu  []string
	Heap []string
	Blocks []string
	GoRoutine []string
	Mutex  []string
	Thread []string
}

func ProfileFromProfileMeta(meta profile.Meta) Profile {
	return Profile{
		ProfileID:  meta.ProfileID,
		ExternalID: meta.ExternalID,
		Type:       meta.Type.String(),
		Service:    meta.Service,
		Labels:     meta.Labels,
		CreatedAt:  meta.CreatedAt.Truncate(time.Second),
	}
}
func ProfileFromProfileData(meta profile.Profile) Profile {
	return Profile{
		ProfileID:  meta.ProfileID,
		ExternalID: meta.ExternalID,
		Type:       meta.Type,
		Service:    meta.Service,
		Labels:     meta.Labels,
		CreatedAt:  meta.CreatedAt.Truncate(time.Second),
	}
}
