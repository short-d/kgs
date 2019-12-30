package db

import (
	"database/sql"
	"github.com/byliuyang/kgs/app/usecase/transactional"
)

type txSQLFn func(tx *sql.Tx) error

func withTransactionSQL(db *sql.DB, fn txSQLFn) (err error) {
	sqlTx, err := db.Begin()

	if err != nil {
		return err
	}

	return transactional.WithTransaction(sqlTx, func(tx transactional.Transaction) error {
		return fn(sqlTx)
	})
}
