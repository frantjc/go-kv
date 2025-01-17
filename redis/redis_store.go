package redis

import (
	"context"
	"fmt"
	"net/url"

	"github.com/frantjc/go-kv"
	xslice "github.com/frantjc/x/slice"
	"github.com/redis/go-redis/v9"
)

type Store struct {
	client *redis.Client
	codec  kv.Codec
}

func (s *Store) Set(ctx context.Context, k string, v any) error {
	data, err := s.codec.Marshal(v)
	if err != nil {
		return err
	}

	if err = s.client.Set(ctx, k, string(data), 0).Err(); err != nil {
		return err
	}

	return nil
}

func (s *Store) Get(ctx context.Context, k string, v any) (found bool, err error) {
	dataString, err := s.client.Get(ctx, k).Result()
	if err != nil {
		if err == redis.Nil {
			return false, nil
		}
		return false, err
	}

	return true, s.codec.Unmarshal([]byte(dataString), v)
}

func (s *Store) Delete(ctx context.Context, k string) error {
	_, err := s.client.Del(ctx, k).Result()
	return err
}

func (s *Store) Close() error {
	return s.client.Close()
}

type Opt func(*Store)

func WithCodec(codec kv.Codec) Opt {
	return func(s *Store) {
		s.codec = codec
	}
}

func New(client *redis.Client, opts ...Opt) (*Store, error) {
	store := &Store{
		client: client,
		codec:  kv.DefaultCodec,
	}

	for _, opt := range opts {
		opt(store)
	}

	return store, nil
}

const Scheme = "redis"

type Opener struct{}

func (o *Opener) Open(ctx context.Context, u *url.URL) (kv.Store, error) {
	if u.Scheme != Scheme {
		return nil, fmt.Errorf("invalid scheme %s, expected %s", u.Scheme, Scheme)
	}

	password, ok := u.User.Password()
	if !ok {
		password = u.Query().Get("password")
	}

	if u.Host == "" {
		u.Host = "localhost:6379"
	} else if u.Port() == "" {
		u.Host += ":6379"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     u.Host,
		Username: xslice.Coalesce(u.User.Username(), u.Query().Get("username")),
		Password: password,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return New(client)
}

func init() {
	kv.Register(&Opener{}, Scheme)
}
