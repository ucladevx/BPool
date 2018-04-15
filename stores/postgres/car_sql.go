package postgres

const (
	carsCreateTable = `
	CREATE TABLE IF NOT EXISTS cars (
		id varchar(20) primary key,
		make varchar(128) NOT NULL,
		model varchar(128),
		year int NOT NULL ,
		color varchar(128) NOT NULL,
		user_id varchar(20) NOT NULL,
		created_at timestamptz DEFAULT NOW(),
		updated_at timestamptz DEFAULT NOW()
	);`

	carsGetAllSQL = "SELECT * FROM cars WHERE id >= $1 LIMIT $2"

	carsGetByIDSQL = "SELECT * FROM cars WHERE id=$1"

	carsInsertSQL = `
		INSERT INTO cars (make, model, year, color, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	carsDeleteSQL = "DELETE FROM cars WHERE id=$1"
)
