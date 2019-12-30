package transactionaltest

import (
	"github.com/byliuyang/kgs/app/usecase/transactional"
)

var _ transactional.Factory = (*FactoryFake)(nil)

type FactoryFake struct{}

func (f FactoryFake) Create() (transactional.Transaction, error) {
	return NewTransactionFake(), nil
}

func NewFactoryFake() FactoryFake {
	return FactoryFake{}
}
