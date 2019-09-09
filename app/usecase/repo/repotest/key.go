package repotest

import (
	"errors"
	"fmt"

	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/repo"
)

var _ repo.AvailableKey = (*AvailableKeyFake)(nil)

type AvailableKeyFake struct {
	keys map[entity.Key]bool
}

func (k *AvailableKeyFake) Create(key entity.Key) (entity.Key, error) {
	if _, ok := k.keys[key]; ok {
		return "", errors.New(fmt.Sprintf("key exists: %s", string(key)))
	}
	k.keys[key] = true
	return key, nil
}

func (k AvailableKeyFake) GetKeys() []entity.Key {
	var keys []entity.Key
	for key := range k.keys {
		keys = append(keys, key)
	}
	return keys
}

func NewAvailableKeyFake() AvailableKeyFake {
	return AvailableKeyFake{
		keys: make(map[entity.Key]bool),
	}
}
