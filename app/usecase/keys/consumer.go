package keys

import (
	"sync"

	"github.com/byliuyang/kgs/app/usecase/repo"
)

type Consumer interface {
	ConsumeInBatch(maxCount uint) ([]string, error)
}

var _ Consumer = (*ConsumerPersist)(nil)

type ConsumerPersist struct {
	mutex            *sync.Mutex
	availableKeyRepo repo.AvailableKey
	allocatedKeyRepo repo.AllocatedKey
}

func (p ConsumerPersist) ConsumeInBatch(maxCount uint) ([]string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	keys, err := p.availableKeyRepo.RetrieveInBatch(maxCount)
	if err != nil {
		return nil, err
	}

	err = p.allocatedKeyRepo.CreateInBatch(keys)
	if err != nil {
		return nil, err
	}

	err = p.availableKeyRepo.DeleteInBatch(keys)
	if err != nil {
		return nil, err
	}

	rawKeys := make([]string, 0)
	for _, key := range keys {
		rawKeys = append(rawKeys, string(key))
	}
	return rawKeys, nil
}

func NewConsumerPersist(
	availableKeyRepo repo.AvailableKey,
	allocatedKeyRepo repo.AllocatedKey,
) ConsumerPersist {
	return ConsumerPersist{
		mutex:            &sync.Mutex{},
		availableKeyRepo: availableKeyRepo,
		allocatedKeyRepo: allocatedKeyRepo,
	}
}
