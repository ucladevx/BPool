package postgres

const (
	userCreateTable = `
CREATE TABLE IF NOT EXISTS users (
    id serial primary key,
    name varchar(128) NOT NULL,
    email varchar(512) NOT NULL ,
    image varchar(1024),
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE ("email")
);`

	userGetAllSQL = "SELECT * FROM users WHERE id >= $1 LIMIT $2"

	userGetByIDSQL = "SELECT * FROM users WHERE id=$1"

	userGetByEmailSQL = "SELECT * FROM users WHERE email=$1"

	userInsertSQL = "INSERT INTO uses (name, email, image) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at"
)
