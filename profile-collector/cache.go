package profefe

import (
	"collector/profile"
	"sync"
)

type cache struct {
	mu       sync.Mutex
	dashboard *CompleteProfileDashBoard
}
func newCache() *cache {
	c := &cache {
		dashboard: &CompleteProfileDashBoard{},
	}
	c.dashboard.Services  = make(map[string]map[string][]string)
	return c
}

func (c *cache )PutProfilesIds(service string,podname string, profileId profile.ID) {
	c.mu.Lock()
	_,ok := c.dashboard.Services[service]
	if !ok {
		c.dashboard.Services[service] = make(map[string][]string)
	}
	c.dashboard.Services[service][podname] = append(c.dashboard.Services[service][podname],string(profileId))

	c.mu.Unlock()
}

func (c *cache) GetProfileIds() (*CompleteProfileDashBoard,error) {

	return c.dashboard,nil
}


