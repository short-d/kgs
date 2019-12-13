package provider

import (
	"github.com/byliuyang/kgs/app/usecase/keys"
)

// CacheSize specifies the size of the local cache for fetched keys
type CacheSize int

// NewConsumer creates a buffered cached keys Consumer
func NewConsumer(bufferSize CacheSize, delegate keys.ConsumerPersist) (keys.ConsumerCached, error) {
	return keys.NewCachedConsumer(int(bufferSize), delegate)
}
