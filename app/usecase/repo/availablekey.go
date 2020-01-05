package repo

import "github.com/short-d/kgs/app/entity"

// AvailableKey represents repository persisting available unique keys
type AvailableKey interface {
	Create(key entity.Key) error
	RetrieveInBatch(maxCount uint) ([]entity.Key, error)
	DeleteInBatch(keys []entity.Key) error
}
