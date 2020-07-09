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
	PodId	   string	 	  `json:"podId"`
	CreatedAt  time.Time      `json:"created_at,omitempty"`
}

type CompleteProfileDashBoard struct {
	Services  map[string]map[string]map[string][]string
}


func ProfileFromProfileMeta(meta profile.Meta) Profile {
	return Profile{
		ProfileID:  meta.ProfileID,
		ExternalID: meta.ExternalID,
		Type:       meta.Type.String(),
		Service:    meta.Service,
		Labels:     meta.Labels,
		PodId:      meta.PodId,
		CreatedAt:  meta.CreatedAt.Truncate(time.Second),
	}
}

