package postgres

const (
	carsCreateTable = ""

	carsGetAllSQL = "SELECT * FROM cars WHERE id >= $1 LIMIT $2"

	carsGetByFieldsSQL = "SELECT * FROM cars WHERE" // customize search by functionality

	carsInsertSQL = `
		INSERT INTO cars (make, model, year, color, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	carsDeleteSQL = "DELETE FROM cars WHERE id=$1"
)
