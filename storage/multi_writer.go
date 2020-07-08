package storage

import (
	"bytes"
	"collector/profile"
	"context"
	"io"
	"io/ioutil"
)

type MultiWriter struct {
	writers []Writer
}

var _ Writer = (*MultiWriter)(nil)

func NewMultiWriter(writers ...Writer) *MultiWriter {
	if len(writers) == 0 {
		panic("storage multiwriter with zero writer")
	}
	return &MultiWriter{
		writers: writers,
	}
}

func (mw *MultiWriter) WriteProfile(ctx context.Context, params *WriteProfileParams, r io.Reader) (profile.Meta, error) {
	// fast path for a case of a single writer in the chain
	if len(mw.writers) == 1 {
		return mw.writers[0].WriteProfile(ctx, params, r)
	}

	var rs io.ReadSeeker
	if iors, ok := r.(io.ReadSeeker); ok {
		rs = iors
	} else {
		// XXX(narqo): a slow path for a (hypothetical) case where incoming reader doesn't implement io.Seeker;
		// in the current implementation of profefe.Collector, the incoming data is already wrapped with bytes.Reader
		// thus this shouldn't happen for the most of incoming requests
		data, err := ioutil.ReadAll(r)
		if err != nil {
			return profile.Meta{}, err
		}
		rs = bytes.NewReader(data)
	}

	meta, err := mw.writers[0].WriteProfile(ctx, params, rs)
	if err != nil {
		return profile.Meta{}, err
	}

	// copy the params and set external id for forwarders
	p := *params
	p.ExternalID = meta.ProfileID

	for _, w := range mw.writers[1:] {
		if _, err := rs.Seek(0, io.SeekStart); err != nil {
			return profile.Meta{}, err
		}
		// the rest of writes act as forwarders, i.e. the returned meta data is skipped
		if _, err := w.WriteProfile(ctx, &p, rs); err != nil {
			return profile.Meta{}, err
		}
	}

	return meta, nil
}

