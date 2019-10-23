package keys

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/keys/gen"
	"github.com/byliuyang/kgs/app/usecase/repo"
)

type Producer interface {
	Produce(KeyLength uint) error
}

var _ Producer = (*ProducerPersist)(nil)

// Producer represents persistent key producer
type ProducerPersist struct {
	logger fw.Logger
	repo   repo.AvailableKey
	keyGen gen.Generator
}

// Produce generates unique keys and store them in the repository
func (p ProducerPersist) Produce(KeyLength uint) error {
	keys := make(chan entity.Key)

	p.keyGen.GenerateKeys(KeyLength, keys)
	for key := range keys {
		p.logger.Info(string(key))
		err := p.repo.Create(key)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewProducer creates and initializes Producer
func NewProducerPersist(
	repo repo.AvailableKey,
	keyGen gen.Generator,
	logger fw.Logger,
) ProducerPersist {
	return ProducerPersist{
		repo:   repo,
		keyGen: keyGen,
		logger: logger,
	}
}