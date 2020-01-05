package repo

import "github.com/short-d/kgs/app/entity"

// AllocatedKey represents repository persisting used keys
type AllocatedKey interface {
	CreateInBatch(keys []entity.Key) error
}
