package profile_collector

import (
	"collector/profile"
	"collector/storage"
	"context"
	"errors"
	"io"
	"log"
	"time"
)

type Collector struct {
	logger *log.Logger
	sw     storage.Writer
}


func NewCollector(logger *log.Logger, sw storage.Writer) *Collector {
	return &Collector{
		logger: logger,
		sw:     sw,
	}

}
func (c *Collector) WriteProfile(ctx context.Context, params *storage.WriteProfileParams, r io.Reader) (profile.Profile, error) {
	return profile.Profile{}, errors.New("method not implemented ")
}

func (c *Collector) writeProfile(ctx context.Context, params *storage.WriteProfileParams, r io.Reader) (profile.Profile, error) {
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	meta, err := c.sw.WriteProfile(ctx, params, r)
	if err != nil {
		return profile.Profile{}, err
	}
	return profile.ProfileFromProfileMeta(meta), nil
}
