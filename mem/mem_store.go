package mem

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/frantjc/go-kv"
)

type Store struct {
	m     *sync.Map
	codec kv.Codec
}

func (s *Store) Set(_ context.Context, k string, v any) error {
	data, err := s.codec.Marshal(v)
	if err != nil {
		return err
	}

	s.m.Store(k, data)
	return nil
}
func (s *Store) Get(_ context.Context, k string, v any) (found bool, err error) {
	dataInterface, found := s.m.Load(k)
	if !found {
		return false, nil
	}

	data := dataInterface.([]byte)

	return true, s.codec.Unmarshal(data, v)
}

func (s *Store) Delete(_ context.Context, k string) error {
	s.m.Delete(k)
	return nil
}

func (s *Store) Close() error {
	s.m = nil
	return nil
}

type Opt func(*Store)

func New(opts ...Opt) *Store {
	store := &Store{
		m:     &sync.Map{},
		codec: kv.DefaultCodec,
	}

	for _, opt := range opts {
		opt(store)
	}

	return store
}

const Scheme = "mem"

type Opener struct{}

func (o *Opener) Open(_ context.Context, u *url.URL) (kv.Store, error) {
	if u.Scheme != Scheme {
		return nil, fmt.Errorf("invalid scheme %s, expected %s", u.Scheme, Scheme)
	}

	return New(), nil
}

func init() {
	kv.Register(&Opener{}, Scheme)
}
