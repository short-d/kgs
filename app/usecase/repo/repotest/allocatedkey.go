package repotest

import (
	"fmt"

	"github.com/short-d/kgs/app/entity"
	"github.com/short-d/kgs/app/usecase/repo"
)

var _ repo.AllocatedKey = (*AllocatedKeyFake)(nil)

type AllocatedKeyFake struct {
	keys map[entity.Key]struct{}
	err  error
}

func (a AllocatedKeyFake) CreateInBatch(keys []entity.Key) error {
	for _, key := range keys {
		if err := a.create(key); err != nil {
			return err
		}
	}

	return a.err
}

func (a *AllocatedKeyFake) FakeError(err error) {
	a.err = err
}

func (a AllocatedKeyFake) create(key entity.Key) error {
	if _, ok := a.keys[key]; ok {
		return fmt.Errorf("key exists: %s", string(key))
	}

	a.keys[key] = struct{}{}

	return nil
}

func NewAllocatedKeyFake() AllocatedKeyFake {
	return AllocatedKeyFake{
		keys: make(map[entity.Key]struct{}),
	}
}
