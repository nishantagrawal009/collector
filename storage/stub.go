package storage

import (
	"collector/profile"
	"context"
	"io"
)

type WriteProfileFunc func(ctx context.Context, params *WriteProfileParams, r io.Reader) (profile.Meta, error)

type StubWriter struct {
	WriteProfileFunc
}

var _ Writer = (*StubWriter)(nil)

func (sw *StubWriter) WriteProfile(ctx context.Context, params *WriteProfileParams, r io.Reader) (profile.Meta, error) {
	return sw.WriteProfileFunc(ctx, params, r)
}

type ListServicesFunc func(ctx context.Context) ([]string, error)

type FindProfilesFunc func(ctx context.Context, params *FindProfilesParams) ([]profile.Meta, error)

type FindProfileIDsFunc func(ctx context.Context, params *FindProfilesParams) ([]profile.ID, error)

type ListProfilesFunc func(ctx context.Context, pid []profile.ID) (ProfileList, error)

type StubReader struct {
	ListServicesFunc
	ListProfilesFunc
	FindProfilesFunc
	FindProfileIDsFunc
}

var _ Reader = (*StubReader)(nil)

func (sr *StubReader) ListServices(ctx context.Context) ([]string, error) {
	return sr.ListServicesFunc(ctx)
}

func (sr *StubReader) FindProfiles(ctx context.Context, params *FindProfilesParams) ([]profile.Meta, error) {
	return sr.FindProfilesFunc(ctx, params)
}

func (sr *StubReader) FindProfileIDs(ctx context.Context, params *FindProfilesParams) ([]profile.ID, error) {
	return sr.FindProfileIDsFunc(ctx, params)
}

func (sr *StubReader) ListProfiles(ctx context.Context, pid []profile.ID) (ProfileList, error) {
	return sr.ListProfilesFunc(ctx, pid)
}

