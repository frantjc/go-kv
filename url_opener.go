package kv

import (
	"context"
	"fmt"
	"net/url"
)

type Opener interface {
	Open(context.Context, *url.URL) (Store, error)
}

type URLMux map[string]Opener

func (m URLMux) Open(ctx context.Context, s string) (Store, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	opener, ok := m[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("no opener registered for scheme %s", u.Scheme)
	}

	return opener.Open(ctx, u)
}

var defaultURLMux = URLMux{}

func Register(opener Opener, scheme string, schemes ...string) {
	for _, s := range append(schemes, scheme) {
		defaultURLMux[s] = opener
	}
}

func Open(ctx context.Context, s string) (Store, error) {
	return defaultURLMux.Open(ctx, s)
}
