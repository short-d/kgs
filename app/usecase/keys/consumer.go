package keys

import (
	"github.com/short-d/kgs/app/entity"
	"github.com/short-d/kgs/app/usecase/repo"
	"github.com/short-d/kgs/app/usecase/transactional"
)

type Consumer interface {
	ConsumeInBatch(maxCount uint) ([]string, error)
}

var _ Consumer = (*ConsumerPersist)(nil)

type AvailableKeyRepoFactory func(tx transactional.Transaction) (repo.AvailableKey, error)
type AllocatedKeyRepoFactory func(tx transactional.Transaction) (repo.AllocatedKey, error)

type ConsumerPersist struct {
	availableKeyRepoFactory AvailableKeyRepoFactory
	allocatedKeyRepoFactory AllocatedKeyRepoFactory
	transactionFactory      transactional.Factory
}

func (p ConsumerPersist) ConsumeInBatch(maxCount uint) (rawKeys []string, err error) {
	tx, err := p.transactionFactory.Create()

	if err != nil {
		return nil, err
	}

	var keys []entity.Key

	err = transactional.WithTransaction(tx, func(tx transactional.Transaction) error {
		availableKeyRepo, err := p.availableKeyRepoFactory(tx)

		if err != nil {
			return err
		}

		allocatedKeyRepo, err := p.allocatedKeyRepoFactory(tx)

		if err != nil {
			return err
		}

		keys, err = availableKeyRepo.RetrieveInBatch(maxCount)

		if err != nil {
			return err
		}

		if err = allocatedKeyRepo.CreateInBatch(keys); err != nil {
			return err
		}

		return availableKeyRepo.DeleteInBatch(keys)
	})

	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		rawKeys = append(rawKeys, string(key))
	}

	return rawKeys, nil
}

func NewConsumerPersist(
	availableKeyRepoFactory AvailableKeyRepoFactory,
	allocatedKeyRepoFactory AllocatedKeyRepoFactory,
	transactionFactory transactional.Factory,
) ConsumerPersist {
	return ConsumerPersist{
		availableKeyRepoFactory: availableKeyRepoFactory,
		allocatedKeyRepoFactory: allocatedKeyRepoFactory,
		transactionFactory:      transactionFactory,
	}
}
