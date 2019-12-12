package keys

import (
	"errors"
)

var _ Consumer = (*ConsumerCached)(nil)

type bufferEntry struct {
	key string
	err error
}

type ConsumerCached struct {
	delegate   Consumer
	bufferSize int
	buffer     chan bufferEntry
}

func (p ConsumerCached) ConsumeInBatch(maxCount uint) ([]string, error) {
	keys := make([]string, 0, maxCount)

	for ; maxCount > 0; maxCount-- {
		// there is probability to get the list of keys with the size less than bufferSize
		// in this case we listen to done channel to perform loop break
		done := make(chan bool)

		// we should load new keys as soon as our buffer is empty
		if len(p.buffer) == 0 {
			go func() {
				p.loadKeys()
				done <- true
			}()
		}

		select {
		case entry := <-p.buffer:
			if entry.err != nil {
				return keys, entry.err
			}

			keys = append(keys, entry.key)
		case <-done:
			break
		}

	}

	return keys, nil
}

func (p ConsumerCached) loadKeys() {
	keys, err := p.delegate.ConsumeInBatch(uint(p.bufferSize))

	if err != nil {
		p.buffer <- bufferEntry{
			key: "",
			err: err,
		}
	}

	for _, key := range keys {
		p.buffer <- bufferEntry{
			key: key,
			err: nil,
		}
	}
}

// NewCachedConsumer returns the cached proxy implementation of Consumer interface
func NewCachedConsumer(bufferSize int, delegate Consumer) (ConsumerCached, error) {
	if bufferSize < 1 {
		return ConsumerCached{}, errors.New("buffer size can't be less than 1")
	}

	return ConsumerCached{
		bufferSize: bufferSize,
		buffer:     make(chan bufferEntry, bufferSize),
		delegate:   delegate,
	}, nil
}
