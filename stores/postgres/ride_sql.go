package postgres

const (
	rideTableName = "rides"

	rideCreateTable = `
CREATE EXTENSION IF NOT EXISTS postgis;
CREATE TABLE IF NOT EXISTS rides (
    id varchar(20) primary key,
    driver_id varchar(20) NOT NULL,
    car_id varchar(20) NOT NULL,
    seats int NOT NULL DEFAULT 1,
		start_city varchar(128) NOT NULL,
		end_city varchar(128) NOT NULL,
		start_dest_lat NUMERIC(9,6) NOT NULL,
		start_dest_lon NUMERIC(9,6) NOT NULL,
		end_dest_lat NUMERIC(9,6) NOT NULL,
		end_dest_lon NUMERIC(9,6) NOT NULL,
		price_per_seat INT DEFAULT 15 NOT NULL,
		info TEXT,
		start_date timestamptz NOT NULL,
		created_at timestamptz DEFAULT NOW(),
		updated_at timestamptz DEFAULT NOW(),
		FOREIGN KEY (driver_id) REFERENCES users (id) ON DELETE CASCADE
);`

	// TODO: add foreign key for car_id when that is merged

	rideGetAllSQL = "SELECT * FROM rides WHERE id > $1 LIMIT $2"

	rideGetByIDSQL = "SELECT * FROM rides WHERE id=$1"

	rideInsertSQL = "INSERT INTO rides (id, driver_id, car_id, seats, start_city, end_city, start_dest_lat, start_dest_lon, end_dest_lat, end_dest_lon, price_per_seat, info, start_date)" +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING created_at, updated_at"
)
