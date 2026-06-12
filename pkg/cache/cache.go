package cache

import "context"

// Cache is used as the main implementation for caching libraries
type Cache interface {
	// Get returns the assigned value for a given key
	Get(context.Context, string) (any, error)

	// Set registers a value, regardless if the key as been
	// already set or not.
	Set(context.Context, string, any) (any, error)

	// Incr increments an value for a given key.
	// If the key does not exists, Incr assumes the default
	// value is 0 and will return 1.
	Incr(context.Context, string) (any, error)

	// SetNX check if the key is already defined, and return
	// the result as a boolean value.
	SetNX(context.Context, string, any) (bool, error)
}
