package repo

import "github.com/byliuyang/kgs/app/entity"

// AvailableKey represents repository persisting available unique keys
type AvailableKey interface {
	Create(key entity.Key) error
	RetrieveInBatch(maxCount int) ([]entity.Key, error)
	DeleteInBatch(keys []entity.Key) error
}
