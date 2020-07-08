package profefe

import (
	"collector/log"
	"collector/profile"
	"collector/storage"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"
	"collector/pprofUtil"

)

type Collector struct {
	logger *log.Logger
	sw     storage.Writer
	cache  *cache
}

func NewCollector(logger *log.Logger, sw storage.Writer) *Collector {
	return &Collector{
		logger: logger,
		sw:     sw,
		cache: newCache(),
	}
}

func (c *Collector) WriteProfile(ctx context.Context, params *storage.WriteProfileParams, r io.Reader) (Profile, error) {
	// don't parse or even read trace profiles, pass them directly to an underlying storage.Writer
	if params.Type == profile.TypeTrace {
		return c.writeProfile(ctx, params, r)
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		return Profile{}, err
	}

	parser := pprofUtil.NewProfileParser(data)

	pp, err := parser.ParseProfile()
	if err != nil {
		return Profile{}, fmt.Errorf("could not parse pprof profile: %w", err)
	}
	if pp.TimeNanos > 0 {
		params.CreatedAt = time.Unix(0, pp.TimeNanos).UTC()
	}

	// move reader's reading position to start to allow storage writers to read the data
	parser.Seek(0, io.SeekStart)

	profile, err := c.writeProfile(ctx, params, parser)

	c.cache.PutProfilesIds(profile.Service,profile.Type,profile.ProfileID)

	return profile,err
}

func (c *Collector) writeProfile(ctx context.Context, params *storage.WriteProfileParams, r io.Reader) (Profile, error) {
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	meta, err := c.sw.WriteProfile(ctx, params, r)
	if err != nil {
		return Profile{}, err
	}
	return ProfileFromProfileMeta(meta), nil
}

