package db

import (
	"database/sql"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/repo"
)

var _ repo.AvailableKey = (*AvailableKeySQL)(nil)

// AvailableKeySQL persist available keys in database through SQL.
type AvailableKeySQL struct {
	db *sql.DB
}

// Create inserts a new key into available_key table.
func (a AvailableKeySQL) Create(key entity.Key) (err error) {
	return withTransactionSQL(a.db, func(tx *sql.Tx) error {
		transactional := NewAvailableKeyTransactional(tx)

		return transactional.Create(key)
	})
}

// RetrieveInBatch fetches maxCount of keys from available_key table.
func (a AvailableKeySQL) RetrieveInBatch(maxCount uint) (keys []entity.Key, err error) {
	err = withTransactionSQL(a.db, func(tx *sql.Tx) error {
		transactional := NewAvailableKeyTransactional(tx)
		keys, err = transactional.RetrieveInBatch(maxCount)

		return err
	})

	if err != nil {
		return nil, err
	}

	return
}

// DeleteInBatch fetch maxCount of keys from available_key table.
func (a AvailableKeySQL) DeleteInBatch(keys []entity.Key) error {
	return withTransactionSQL(a.db, func(tx *sql.Tx) error {
		transactional := NewAvailableKeyTransactional(tx)

		return transactional.DeleteInBatch(keys)
	})
}

// NewAvailableKeySQL creates AvailableKeySQL
func NewAvailableKeySQL(db *sql.DB) AvailableKeySQL {
	return AvailableKeySQL{db: db}
}
