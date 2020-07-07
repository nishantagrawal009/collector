package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

)

type Storage interface {
	Writer
	Reader
}

type Writer interface {
	WriteProfile(ctx context.Context, params *WriteProfileParams, r io.Reader) (Meta, error)
}

type WriteProfileParams struct {
	ExternalID ID
	Service    string
	Type       ProfileType
	Labels     Labels
	CreatedAt  time.Time
}
func (params *WriteProfileParams) Validate() error {
	if params == nil {
		return errors.New("empty params")
	}
	if params.Service == "" {
		return errors.New("empty service")
	}
	if params.Type == TypeUnknown {
		return fmt.Errorf("unknown profile type %s", params.Type)
	}
	return nil
}

type Reader interface {
	FindProfiles(ctx context.Context, params *FindProfilesParams) ([]Meta, error)
	FindProfileIDs(ctx context.Context, params *FindProfilesParams) ([]ID, error)
	ListProfiles(ctx context.Context, pid []ID) (ProfileList, error)
	ListServices(ctx context.Context) ([]string, error)
}

type FindProfilesParams struct {
	Service      string
	Type         ProfileType
	Labels       Labels
	CreatedAtMin time.Time
	CreatedAtMax time.Time
	Limit        int
}

func (params *FindProfilesParams) Validate() error {
	if params == nil {
		return errors.New("empty params")
	}
	if params.Service == "" {
		return errors.New("empty service")
	}
	if params.Type == TypeUnknown {
		return fmt.Errorf("unknown profile type %s", params.Type)
	}
	if params.CreatedAtMin.IsZero() || params.CreatedAtMax.IsZero() {
		return fmt.Errorf("created_at is zero: min %v, max %v", params.CreatedAtMin, params.CreatedAtMax)
	}
	if params.CreatedAtMin.After(params.CreatedAtMax) {
		return fmt.Errorf("created_at min after max: min %v, max %v", params.CreatedAtMin, params.CreatedAtMax)
	}
	return nil
}

type ProfileList interface {
	Next() bool
	Profile() (io.Reader, error)
	Close() error
}
