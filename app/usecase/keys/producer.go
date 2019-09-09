package keys

import (
	"fmt"

	"github.com/byliuyang/app/fw"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/keys/gen"
	"github.com/byliuyang/kgs/app/usecase/repo"
)

// Producer represents persistent key producer
type Producer struct {
	logger fw.Logger
	repo   repo.AvailableKey
	keyGen gen.Generator
}

// Produce generates unique keys and store them in the repository
func (p Producer) Produce() {
	keys := make(chan entity.Key)
	p.keyGen.GenerateKeys(keys)

	for key := range keys {
		key, err := p.repo.Create(key)
		if err != nil {
			p.logger.Error(err)
		} else {
			p.logger.Info(fmt.Sprintf("Saved %v", key))
		}
	}
}

// NewProducer creates and initializes Producer
func NewProducer(
	repo repo.AvailableKey,
	keyGen gen.Generator,
	logger fw.Logger,
) Producer {
	return Producer{
		repo:   repo,
		keyGen: keyGen,
		logger: logger,
	}
}
