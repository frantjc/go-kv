package kv

import "context"

type Store interface {
	Set(context.Context, string, any) error
	Get(context.Context, string, any) (bool, error)
	Delete(context.Context, string) error
	Close() error
}
