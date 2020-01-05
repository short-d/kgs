package gen

import "github.com/short-d/kgs/app/entity"

// Generator represents unique key generator
type Generator interface {
	GenerateKeys(keySize uint, keysOut chan<- entity.Key)
}
