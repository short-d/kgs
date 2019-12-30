package db

import (
	"database/sql"
	"fmt"
	"github.com/byliuyang/kgs/app/adapter/db/table"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/repo"
	"time"
)

var _ repo.AllocatedKey = (*AllocatedKeyTransactional)(nil)

// AllocatedKeyTransactional performs SQL operations using the underlying sql transaction object
type AllocatedKeyTransactional struct {
	tx *sql.Tx
}

func (a AllocatedKeyTransactional) CreateInBatch(keys []entity.Key) error {
	statement, err := a.tx.Prepare(
		fmt.Sprintf(`
INSERT INTO "%s" ("%s", "%s")
VALUES ($1, $2);
`,
			table.AllocatedKey.TableName,
			table.AllocatedKey.ColumnKey,
			table.AllocatedKey.ColumnAllocatedAt),
	)
	if err != nil {
		return err
	}
	// Prepared statements take up server resources and should be closed after
	// use.
	defer statement.Close()

	for _, key := range keys {
		now := time.Now()
		_, err := statement.Exec(key, now)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewAllocatedKeyTransactional creates a new instance of AllocatedKeyTransactional
func NewAllocatedKeyTransactional(tx *sql.Tx) AllocatedKeyTransactional {
	return AllocatedKeyTransactional{tx: tx}
}
