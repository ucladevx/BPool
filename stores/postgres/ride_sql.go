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
	price_per_seat NUMERIC(10,2) DEFAULT 15 NOT NULL,
	info TEXT,
	start_date timestamptz NOT NULL,
	created_at timestamptz DEFAULT NOW(),
	updated_at timestamptz DEFAULT NOW(),
	FOREIGN KEY (driver_id) REFERENCES users (id) ON DELETE CASCADE,
	FOREIGN KEY (car_id) REFERENCES cars (id) ON DELETE CASCADE
);`

	rideGetAllSQL = "SELECT * FROM rides WHERE id > $1 LIMIT $2"

	rideGetByIDSQL = "SELECT *, (SELECT COUNT(*) FROM passengers WHERE ride_id=$1 AND status='accepted') AS seats_taken FROM rides WHERE rides.id=$1;"

	rideInsertSQL = "INSERT INTO rides (id, driver_id, car_id, seats, start_city, end_city, start_dest_lat, start_dest_lon, end_dest_lat, end_dest_lon, price_per_seat, info, start_date) " +
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING created_at, updated_at"
	rideUpdateSQL = "UPDATE rides SET car_id=$1, seats=$2, start_city=$3, end_city=$4, start_dest_lat=$5, start_dest_lon=$6, end_dest_lat=$7, end_dest_lon=$8, price_per_seat=$9, info=$10, start_date=$11, updated_at=NOW() " +
		"WHERE id=$12 RETURNING updated_at"

	rideDeleteSQL = "DELETE FROM rides WHERE id=$1"
)
