package keys

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/keys/gen"
	"github.com/byliuyang/kgs/app/usecase/transactional"
)

type Producer interface {
	Produce(KeyLength uint) error
}

var _ Producer = (*ProducerPersist)(nil)

// Producer represents persistent key producer
type ProducerPersist struct {
	availableKeyFactory AvailableKeyRepoFactory
	transactionFactory  transactional.Factory
	logger              fw.Logger
	keyGen              gen.Generator
}

// Produce generates unique keys and store them in the repository
func (p ProducerPersist) Produce(KeyLength uint) error {
	tx, err := p.transactionFactory.Create()

	if err != nil {
		return err
	}

	return transactional.WithTransaction(tx, func(tx transactional.Transaction) error {
		repo, err := p.availableKeyFactory(tx)

		if err != nil {
			return err
		}

		keys := make(chan entity.Key)
		p.keyGen.GenerateKeys(KeyLength, keys)

		for key := range keys {
			if err := repo.Create(key); err != nil {
				return err
			}
		}

		return nil
	})
}

// NewProducer creates and initializes Producer
func NewProducerPersist(
	availableKeyFactory AvailableKeyRepoFactory,
	transactionFactory transactional.Factory,
	keyGen gen.Generator,
	logger fw.Logger,
) ProducerPersist {
	return ProducerPersist{
		availableKeyFactory: availableKeyFactory,
		transactionFactory: transactionFactory,
		keyGen: keyGen,
		logger: logger,
	}
}
