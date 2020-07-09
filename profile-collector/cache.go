package profefe

import (
	"collector/profile"
	"sync"
	"time"
)

type cache struct {
	mu       sync.Mutex
	dashboard *CompleteProfileDashBoard
}
func newCache() *cache {
	c := &cache {
		dashboard: &CompleteProfileDashBoard{},
	}
	c.dashboard.Services  = make(map[string]map[string]map[string][]PprofFileMeta)
	return c
}

func (c *cache )PutProfilesIds(service string,podname string, profileType string, profileId profile.ID,time  time.Time ) {

	c.mu.Lock()

	_,ok := c.dashboard.Services[service]

	if !ok {
		c.dashboard.Services[service] = make(map[string]map[string][]PprofFileMeta)
	}

	_,ok = c.dashboard.Services[service][podname]
	if !ok {
		c.dashboard.Services[service][podname] = make(map[string][]PprofFileMeta)
	}
	pprofMeta := PprofFileMeta{Id:string(profileId),Time:time}
	c.dashboard.Services[service][podname][profileType] = append(c.dashboard.Services[service][podname][profileType],pprofMeta)

	c.mu.Unlock()
}

func (c *cache) GetProfileIds() (*CompleteProfileDashBoard,error) {

	return c.dashboard,nil
}


