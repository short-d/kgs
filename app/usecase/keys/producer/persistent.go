package producer

import (
	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/keys/gen"
	"github.com/byliuyang/kgs/app/usecase/repo"
)

var _ Producer = (*Persist)(nil)

// Producer represents persistent key producer
type Persist struct {
	logger fw.Logger
	repo   repo.AvailableKey
	keyGen gen.Generator
}

// Produce generates unique keys and store them in the repository
func (p Persist) Produce(keySize uint) {
	keys := make(chan entity.Key)

	p.keyGen.GenerateKeys(keySize, keys)

	for key := range keys {
		p.logger.Info(string(key))
		err := p.repo.Create(key)
		if err != nil {
			p.logger.Error(err)
			return
		}
	}
}

// NewProducer creates and initializes Producer
func NewPersist(
	repo repo.AvailableKey,
	keyGen gen.Generator,
	logger fw.Logger,
) Persist {
	return Persist{
		repo:   repo,
		keyGen: keyGen,
		logger: logger,
	}
}
