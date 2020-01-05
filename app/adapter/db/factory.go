package db

import (
	"database/sql"

	"github.com/short-d/kgs/app/usecase/transactional"
)

var _ transactional.Factory = (*FactorySQL)(nil)

// FactorySQL creates a SQL transaction entity
type FactorySQL struct {
	db *sql.DB
}

func (f FactorySQL) Create() (transactional.Transaction, error) {
	return f.db.Begin()
}

// NewFactorySQL creates a new instance of FactorySQL
func NewFactorySQL(db *sql.DB) FactorySQL {
	return FactorySQL{
		db: db,
	}
}
