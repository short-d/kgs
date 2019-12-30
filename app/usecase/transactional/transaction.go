package transactional

// Transaction is a abstract tx
type Transaction interface {
	Rollback() error
	Commit() error
}
