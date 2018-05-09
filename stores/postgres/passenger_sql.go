package postgres

const (
	passengerTableName = "passengers"

	passengerCreateTable = `
CREATE TABLE IF NOT EXISTS passengers (
	id varchar(20) primary key,
	driver_id varchar(20) NOT NULL,
	passenger_id varchar(20) NOT NULL,
	ride_id varchar(20) NOT NULL,
	status varchar(20) NOT NULL,
	created_at timestamptz DEFAULT NOW(),
	updated_at timestamptz DEFAULT NOW(),
	FOREIGN KEY (driver_id) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (passenger_id) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (ride_id) REFERENCES rides (id) ON DELETE CASCADE,
	UNIQUE (driver_id, passenger_id, ride_id)
);`

	passengerGetAllSQL = "SELECT * FROM passengers WHERE id > $1 LIMIT $2"

	passengerGetByIDSQL = "SELECT * FROM passengers WHERE id=$1"

	passengerInsertSQL = "INSERT INTO passengers (id, driver_id, passenger_id, ride_id, status) VALUES ($1, $2, $3, $4, $5) RETURNING created_at, updated_at"
	passengerUpdateSQL = "UPDATE passengers SET status=$1 updated_at=NOW() WHERE id=$2 RETURNING updated_at"

	passengerDeleteSQL = "DELETE FROM passengers WHERE id=$1"
)
