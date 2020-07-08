package profefe

import (
	"collector/profile"
	"fmt"
	"sync"
)

type cache struct {
	mu       sync.Mutex
	profiles map[string]profile.ID
}
func newCache() *cache {
	c := &cache{
		profiles: make(map[string]profile.ID),
	}
	return c
}

func (c *cache )PutProfilesIds(service string, profileId profile.ID) {
	c.mu.Lock()
	c.profiles[service] = profileId
	c.mu.Unlock()
}

func (c *cache) GetProfileIds(service string) (profile.ID,error) {

	c.mu.Lock()
	profiles, ok := c.profiles[service]
	if !ok {
		return "", fmt.Errorf("could not find the service")
	}
	c.mu.Unlock()
	return profiles,nil
}


