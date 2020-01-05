package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/short-d/kgs/app/adapter/db/table"
	"github.com/short-d/kgs/app/entity"
	"github.com/short-d/kgs/app/usecase/repo"
)

var _ repo.AvailableKey = (*AvailableKeySQL)(nil)

// AvailableKeySQL persist available keys in database through SQL.
type AvailableKeySQL struct {
	db *sql.DB
}

// Create inserts a new key into available_key table.
func (a AvailableKeySQL) Create(key entity.Key) error {
	statement := fmt.Sprintf(`
INSERT INTO "%s" ("%s", "%s")
VALUES ($1, $2);
`,
		table.AvailableKey.TableName,
		table.AvailableKey.ColumnKey,
		table.AvailableKey.ColumnCreatedAt,
	)

	now := time.Now()
	_, err := a.db.Exec(statement, key, now)
	return err
}

// RetrieveInBatch fetches maxCount of keys from available_key table.
func (a AvailableKeySQL) RetrieveInBatch(maxCount uint) ([]entity.Key, error) {
	query := fmt.Sprintf(`
SELECT "%s"
FROM "%s"
LIMIT $1
`,
		table.AvailableKey.ColumnKey,
		table.AvailableKey.TableName,
	)
	rows, err := a.db.Query(query, maxCount)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	keys := make([]entity.Key, 0)

	for rows.Next() {
		var key string
		err := rows.Scan(&key)

		if err != nil {
			return nil, err
		}
		keys = append(keys, entity.Key(key))
	}

	err = rows.Close()
	if err != nil {
		return nil, err
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return keys, nil
}

// DeleteInBatch fetch maxCount of keys from available_key table.
func (a AvailableKeySQL) DeleteInBatch(keys []entity.Key) error {
	tx, err := a.db.Begin()
	if err != nil {
		return nil
	}
	// The rollback will be ignored if the tx has been committed later in the
	// function.
	defer tx.Rollback()

	statement, err := tx.Prepare(
		fmt.Sprintf(`
DELETE FROM "%s"
WHERE "%s"=$1;
`,
			table.AvailableKey.TableName,
			table.AvailableKey.ColumnKey),
	)
	if err != nil {
		return err
	}
	// Prepared statements take up server resources and should be closed after
	// use.
	defer statement.Close()

	for _, key := range keys {
		_, err := statement.Exec(key)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// NewAvailableKeySQL creates AvailableKeySQL
func NewAvailableKeySQL(db *sql.DB) AvailableKeySQL {
	return AvailableKeySQL{db: db}
}
