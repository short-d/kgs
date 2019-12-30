package transactional

// Factory represents the interface for creating a transaction object
type Factory interface {
	// Create creates a Transaction instance
	Create() (Transaction, error)
}
