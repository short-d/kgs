package db

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/short-d/kgs/app/adapter/db/table"
	"github.com/short-d/kgs/app/entity"
	"github.com/short-d/kgs/app/usecase/repo"
)

var _ repo.AvailableKey = (*AvailableKeyTransactional)(nil)

// AvailableKeyTransactional persist available keys in database through SQL.
type AvailableKeyTransactional struct {
	tx *sql.Tx
}

// Create inserts a new key into available_key table.
func (a AvailableKeyTransactional) Create(key entity.Key) error {
	statement := fmt.Sprintf(`
INSERT INTO "%s" ("%s", "%s")
VALUES ($1, $2);
`,
		table.AvailableKey.TableName,
		table.AvailableKey.ColumnKey,
		table.AvailableKey.ColumnCreatedAt,
	)

	now := time.Now()
	_, err := a.tx.Exec(statement, key, now)
	return err
}

// RetrieveInBatch fetches maxCount of keys from available_key table.
func (a AvailableKeyTransactional) RetrieveInBatch(maxCount uint) ([]entity.Key, error) {
	query := fmt.Sprintf(`
SELECT "%s"
FROM "%s"
LIMIT $1
`,
		table.AvailableKey.ColumnKey,
		table.AvailableKey.TableName,
	)
	rows, err := a.tx.Query(query, maxCount)
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
func (a AvailableKeyTransactional) DeleteInBatch(keys []entity.Key) error {
	statement, err := a.tx.Prepare(
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
	return nil
}

// NewAvailableKeyTransactional creates a new AvailableKeyTransactional object
func NewAvailableKeyTransactional(tx *sql.Tx) AvailableKeyTransactional {
	return AvailableKeyTransactional{tx: tx}
}
