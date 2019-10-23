package table

// AllocatedKey represents allocated_key table in the database
var AllocatedKey = struct {
	TableName         string
	ColumnKey         string
	ColumnAllocatedAt string
}{
	TableName:         "allocated_key",
	ColumnKey:         "key",
	ColumnAllocatedAt: "allocated_at",
}
