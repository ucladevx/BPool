package postgres

// Store is a postgres store
type Store interface {
	migrate()
}

// CreateTables creates tables for all stores
func CreateTables(stores ...Store) {
	for _, s := range stores {
		s.migrate()
	}
}
