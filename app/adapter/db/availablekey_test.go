package db

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"testing"

	"github.com/byliuyang/app/mdtest"
	"github.com/byliuyang/kgs/app/adapter/db/table"
	"github.com/byliuyang/kgs/app/entity"
)

func TestAvailableKeySQL_Create(t *testing.T) {
	testCases := []struct {
		name     string
		key      entity.Key
		rowExist bool
		hasErr   bool
	}{
		{
			name:     "key exists",
			key:      entity.Key("testKey"),
			rowExist: true,
			hasErr:   true,
		},
		{
			name:     "no given key",
			key:      entity.Key("testKey"),
			rowExist: false,
			hasErr:   false,
		},
	}

	for _, testCase := range testCases {
		db, stub, err := mdtest.NewSQLStub()
		mdtest.Equal(t, nil, err)
		defer db.Close()

		expStatement := fmt.Sprintf(
			`INSERT\s*INTO\s*"%s"`,
			table.AvailableKey.TableName,
		)
		if testCase.rowExist {
			stub.ExpectExec(expStatement).WillReturnError(errors.New("row exists"))
		} else {
			stub.ExpectExec(expStatement).WillReturnResult(driver.ResultNoRows)
		}

		userRepo := NewAvailableKeySQL(db)

		err = userRepo.Create(testCase.key)
		if testCase.hasErr {
			mdtest.NotEqual(t, nil, err)
			return
		}
		mdtest.Equal(t, nil, err)
	}
}

//func TestAvailableKeySQL_RetrieveInBatch(t *testing.T) {
//	testCases := []struct {
//		name      string
//		tableRows *mdtest.TableRows
//		maxCount  int
//		hasErr    bool
//		expKeys   []entity.Key
//	}{
//		{
//			name:  "maxCount = 0",
//			maxCount: 0,
//			tableRows: mdtest.NewTableRows([]string{
//				table.AvailableKey.ColumnKey,
//				table.AvailableKey.ColumnCreatedAt,
//			}).
//				AddRow("testKey", nil),
//			hasErr: false,
//			expKeys: []entity.Key{},
//		},
//		{
//			name:  "table is empty",
//			maxCount: 10,
//			tableRows: mdtest.NewTableRows([]string{
//				table.AvailableKey.ColumnKey,
//				table.AvailableKey.ColumnCreatedAt,
//			}),
//			hasErr: false,
//			expKeys: []entity.Key{},
//		},
//		{
//			name:  "table has less keys than maxCount",
//			maxCount: 10,
//			tableRows: mdtest.NewTableRows([]string{
//				table.AvailableKey.ColumnKey,
//				table.AvailableKey.ColumnCreatedAt,
//			}).
//				AddRow("firstKey", nil).
//				AddRow("secondKey", nil).
//				AddRow("thirdKey", nil),
//			hasErr: false,
//			expKeys: []entity.Key{
//				entity.Key("firstKey"),
//				entity.Key("secondKey"),
//				entity.Key("thirdKey"),
//			},
//		},
//		{
//			name:  "table has more keys than maxCount",
//			maxCount: 2,
//			tableRows: mdtest.NewTableRows([]string{
//				table.AvailableKey.ColumnKey,
//				table.AvailableKey.ColumnCreatedAt,
//			}).
//				AddRow("firstKey", nil).
//				AddRow("secondKey", nil).
//				AddRow("thirdKey", nil),
//			hasErr: false,
//			expKeys: []entity.Key{
//				entity.Key("firstKey"),
//				entity.Key("secondKey"),
//			},
//		},
//	}
//
//	for _, testCase := range testCases {
//		t.Run(testCase.name, func(t *testing.T) {
//			db, stub, err := mdtest.NewSQLStub()
//			mdtest.Equal(t, nil, err)
//			defer db.Close()
//
//			expQuery := fmt.Sprintf(
//				`^SELECT ".+",".+",".+",".+",".+" FROM "%s" LIMIT .+$`,
//				table.AvailableKey.TableName,
//				)
//
//			rowCount := math.Min(
//				float64(testCase.maxCount),
//				float64(len(testCase.tableRows)))
//			stub.ExpectQuery(expQuery).WillReturnRows(testCase.tableRows)
//
//			userRepo := NewUserSQL(db)
//
//			gotUser, err := userRepo.GetByEmail(testCase.email)
//			if testCase.hasErr {
//				mdtest.NotEqual(t, nil, err)
//				return
//			}
//			mdtest.Equal(t, nil, err)
//			mdtest.Equal(t, testCase.expUser, gotUser)
//		})
//	}
//}
