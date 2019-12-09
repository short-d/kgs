package provider

import (
	"github.com/byliuyang/kgs/app/usecase/keys"
)

// KeyFetchBufferSize specifies the size of the local cache for fetched keys
type KeyFetchBufferSize int

// NewConsumer creates a buffered cached keys Consumer
func NewConsumer(bufferSize KeyFetchBufferSize, delegate keys.ConsumerPersist) (keys.ConsumerCached, error) {
	return keys.NewCachedConsumer(int(bufferSize), delegate)
}
