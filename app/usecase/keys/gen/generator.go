package gen

import "github.com/byliuyang/kgs/app/entity"

// Generator represents unique key generator
type Generator interface {
	GenerateKeys(keysOut chan<- entity.Key)
}
