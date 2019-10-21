package gen

import "github.com/byliuyang/kgs/app/entity"

// Generator represents unique key generator
type Generator interface {
	GenerateKeys(keySize uint, keysOut chan<- entity.Key)
}
