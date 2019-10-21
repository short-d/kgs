package table

// AvailableKey represents available_key database table
var AvailableKey = struct {
	TableName       string
	ColumnKey       string
	ColumnCreatedAt string
}{
	TableName:       "available_key",
	ColumnKey:       "key",
	ColumnCreatedAt: "created_at",
}
