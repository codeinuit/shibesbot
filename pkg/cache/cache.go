package cache

import "context"

// Cache is an interface for cache implementations
type Cache interface {
	Get(context.Context, string) (any, error)
	Set(context.Context, string, string) (any, error)
	Incr(context.Context, string) (any, error)
	SetNX(context.Context, string, any) (any, error)
}
