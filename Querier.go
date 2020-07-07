package main

import (
	"context"
	"errors"
	"io"
)

type Querier struct {
	sr     Reader
}

func NewQuerier(sr Reader) *Querier {
	return &Querier{
		sr:     sr,
	}
}


func (q *Querier) GetProfilesTo(ctx context.Context, dst io.Writer, pids []ID) error {
return errors.New("method not implemented")
}

func (q *Querier) FindMergeProfileTo(ctx context.Context, dst io.Writer, params *FindProfilesParams) error {
	return errors.New("method not implemented")
}

func (q *Querier) ListServices(ctx context.Context) ([]string, error) {
	return nil,errors.New("method not implemented")
}
