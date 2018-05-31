package postgres

const (
	ratingCreateTable = `
		CREATE TABLE IF NOT EXISTS ratings (
			id varchar(20) primary key,
			rating int NOT NULL,
			ride_id varchar(20) NOT NULL,
			rater_id varchar(20) NOT NULL,
			ratee_id varchar(20) NOT NULL,
			comment varchar(500),
			created_at timestamptz DEFAULT NOW(),
			updated_at timestamptz DEFAULT NOW(),
			FOREIGN KEY (ride_id) REFERENCES rides (id) ON DELETE CASCADE,
			FOREIGN KEY (rater_id) REFERENCES users (id) ON DELETE CASCADE,
			FOREIGN KEY (ratee_id) REFERENCES users (id) ON DELETE CASCADE
		)
	`

	ratingInsertSQL = `
		INSERT INTO ratings (id, rating, ride_id, rater_id, ratee_id, comment)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING created_at, updated_at
	`

	ratingGetByIDSQL = `
		SELECT * FROM ratings WHERE id=$1
	`
	ratingGetAverageRatingSQL = `
		SELECT AVG(rating) FROM ratings WHERE ratee_id=$1
	`

	ratingDeleteSQL = `
		DELETE FROM ratings WHERE id=$1
	`
)
