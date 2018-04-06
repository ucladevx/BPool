package postgres

const (
	userTableName = "users"

	userCreateTable = `
CREATE TABLE IF NOT EXISTS users (
    id varchar(20) primary key,
    first_name varchar(128) NOT NULL,
    last_name varchar(128) NOT NULL,
    email varchar(512) NOT NULL ,
    profile_image varchar(1024),
    auth_level integer DEFAULT 0,
    created_at timestamptz DEFAULT NOW(),
    updated_at timestamptz DEFAULT NOW(),
    UNIQUE ("email")
);`

	userGetAllSQL = "SELECT * FROM users WHERE id >= $1 LIMIT $2"

	userGetByIDSQL = "SELECT * FROM users WHERE id=$1"

	userGetByEmailSQL = "SELECT * FROM users WHERE email=$1"

	userInsertSQL = "INSERT INTO users (id, first_name, last_name, email, profile_image) VALUES ($1, $2, $3, $4, $5) RETURNING auth_level, created_at, updated_at"
)
