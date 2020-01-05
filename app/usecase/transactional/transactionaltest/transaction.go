package transactionaltest

import "github.com/short-d/kgs/app/usecase/transactional"

var _ transactional.Transaction = (*TransactionFake)(nil)

type TransactionFake struct{}

func (t TransactionFake) Commit() error {
	return nil
}

func (t TransactionFake) Rollback() error {
	return nil
}

func NewTransactionFake() TransactionFake {
	return TransactionFake{}
}
