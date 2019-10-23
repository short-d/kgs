package db

import (
	"database/sql"
	"fmt"
	"github.com/byliuyang/kgs/app/adapter/db/table"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/repo"
	"time"
)

var _ repo.AllocatedKey = (*AllocatedKeySQL)(nil)

type AllocatedKeySQL struct {
	db *sql.DB
}

func (a AllocatedKeySQL) CreateInBatch(keys []entity.Key) error {
	tx, err := a.db.Begin()
	if err != nil {
		return nil
	}
	// The rollback will be ignored if the tx has been committed later in the
	// function.
	defer tx.Rollback()

	statement, err := tx.Prepare(
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
	return tx.Commit()
}

func NewAllocatedKeySQL(db *sql.DB) AllocatedKeySQL {
	return AllocatedKeySQL{db: db}
}
