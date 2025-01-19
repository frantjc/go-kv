package kv

import (
	"context"
	"fmt"
	"net/url"
)

type URLOpener interface {
	Open(context.Context, *url.URL) (Store, error)
}

var (
	urlMux = map[string]URLOpener{}
)

func Register(opener URLOpener, scheme string, schemes ...string) {
	for _, s := range append(schemes, scheme) {
		urlMux[s] = opener
	}
}

func Open(ctx context.Context, s string) (Store, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	opener, ok := urlMux[u.Scheme]
	if !ok {
		return nil, fmt.Errorf("no opener registered for scheme %s", u.Scheme)
	}

	return opener.Open(ctx, u)
}
