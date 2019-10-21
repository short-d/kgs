package repo

import "github.com/byliuyang/kgs/app/entity"

// AllocatedKey represents repository persisting used keys
type AllocatedKey interface {
	CreateInBatch(keys []entity.Key) error
	RetrieveInBatch(maxCount int) []entity.Key
	DeleteInBatch(keys []entity.Key) error
}
