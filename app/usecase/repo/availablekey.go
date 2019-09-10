package repo

import "github.com/byliuyang/kgs/app/entity"

// Key represents repository persisting unique keys
type AvailableKey interface {
	Create(key entity.Key) (entity.Key, error)
}
