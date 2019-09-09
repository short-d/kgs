package gentest

import (
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/keys/gen"
)

var _ gen.Generator = (*GeneratorStub)(nil)

type GeneratorStub struct {
	keys []entity.Key
}

func (g GeneratorStub) GenerateKeys(keysOut chan<- entity.Key) {
	go func() {
		for _, key := range g.keys {
			keysOut <- key
		}
		close(keysOut)
	}()
}

func NewGeneratorStub(keys []entity.Key) GeneratorStub {
	return GeneratorStub{
		keys: keys,
	}
}
