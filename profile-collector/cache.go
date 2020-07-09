package profefe

import (
	"collector/profile"
	"sync"
)

type cache struct {
	mu       sync.Mutex
	profiles map[string][]string
}
func newCache() *cache {
	c := &cache{
		profiles: make(map[string][]string),
	}
	return c
}

func (c *cache )PutProfilesIds(service string, profiletype string, profileId profile.ID) {
	c.mu.Lock()
	c.profiles[profiletype] = append(c.profiles[profiletype],string(profileId))
	c.mu.Unlock()
}

func (c *cache) GetProfileIds() (GetProfileDisplay,error) {
	gpd := GetProfileDisplay{}
	c.mu.Lock()
	gpd.Cpu  = c.profiles["cpu"]
	gpd.Heap = c.profiles["heap"]
	gpd.Blocks = c.profiles["block"]
	gpd.GoRoutine = c.profiles["goroutine"]
	gpd.Mutex = c.profiles["mutex"]
	gpd.Thread = c.profiles["threadcreate"]
	c.mu.Unlock()
	return gpd,nil
}


