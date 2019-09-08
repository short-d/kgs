package keygen

import "github.com/byliuyang/kgs/app/entity"

type KeyGenerator interface {
	AvailableKeys(keys chan<- entity.Key)
}
