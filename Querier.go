package main

import (
	"context"
	"io"

	"github.com/pkg/mod/github.com/pkg/errors@v0.9.1"
	"github.com/profefe/profefe/pkg/profile"
	"github.com/profefe/profefe/pkg/storage"
)

type Querier struct {
	sr     Reader
}

func NewQuerier(sr Reader) *Querier {
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
