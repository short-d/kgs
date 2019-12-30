package transactional

type TxFn func(tx Transaction) error

// WithTransaction performs the given callback in a transactional way
func WithTransaction(tx Transaction, fn TxFn) (err error) {
	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and re-panic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	err = fn(tx)

	return err
}
