package repotest

import (
	"fmt"
	"sort"

	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/repo"
)

var _ repo.AvailableKey = (*AvailableKeyFake)(nil)

type AvailableKeyFake struct {
	keys map[entity.Key]struct{}
}

func (a AvailableKeyFake) Create(key entity.Key) error {
	if _, ok := a.keys[key]; ok {
		return fmt.Errorf("key exists: %s", string(key))
	}
	a.keys[key] = struct{}{}
	return nil
}

func (a AvailableKeyFake) RetrieveInBatch(maxCount uint) ([]entity.Key, error) {
	keys := a.GetKeys()
	if len(keys) <= int(maxCount) {
		return keys, nil
	}
	return keys[:maxCount], nil
}

func (a AvailableKeyFake) DeleteInBatch(keys []entity.Key) error {
	for _, key := range keys {
		delete(a.keys, key)
	}
	return nil
}

func (k AvailableKeyFake) GetKeys() []entity.Key {
	var keys []entity.Key
	for key := range k.keys {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	return keys
}

func NewAvailableKeyFake() AvailableKeyFake {
	return AvailableKeyFake{
		keys: make(map[entity.Key]struct{}),
	}
}
