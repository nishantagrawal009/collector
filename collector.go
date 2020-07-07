package main

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

)

type Collector struct {
	logger *log.Logger
	sw     Writer
}


func NewCollector(logger *log.Logger, sw Writer) *Collector {
	return &Collector{
		logger: logger,
		sw:     sw,
	}
}
func (c *Collector) WriteProfile(ctx context.Context, params *WriteProfileParams, r io.Reader) (Profile, error) {
	return Profile{}, errors.New("method not implemented ")
}

func (c *Collector) writeProfile(ctx context.Context, params *WriteProfileParams, r io.Reader) (Profile, error) {
	if params.CreatedAt.IsZero() {
		params.CreatedAt = time.Now().UTC()
	}

	meta, err := c.sw.WriteProfile(ctx, params, r)
	if err != nil {
		return Profile{}, err
	}
	return ProfileFromProfileMeta(meta), nil
}
