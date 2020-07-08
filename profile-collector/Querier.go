package profile_collector

import (
	"collector/profile"
	"collector/storage"
	"context"
	"errors"
	"io"
)

type Querier struct {
	sr storage.Reader
}

func NewQuerier(sr storage.Reader) *Querier {
	return &Querier{
		sr:     sr,
	}
}


func (q *Querier) GetProfilesTo(ctx context.Context, dst io.Writer, pids []profile.ID) error {
return errors.New("method not implemented")
}

func (q *Querier) FindMergeProfileTo(ctx context.Context, dst io.Writer, params *storage.FindProfilesParams) error {
	return errors.New("method not implemented")
}

func (q *Querier) ListServices(ctx context.Context) ([]string, error) {
	return nil,errors.New("method not implemented")
}
