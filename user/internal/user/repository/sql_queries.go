package repository

const (
	createUserQuery = `INSERT INTO users (first_name, last_name, email, password, avatar, role) 
	VALUES ($1,$2,$3,$4,$5,$6) 
	RETURNING user_id, first_name, last_name, email, avatar, role, updated_at, created_at`

	getUserByIDQuery = `SELECT user_id, first_name, last_name, email, avatar, role, updated_at, created_at FROM users WHERE user_id = $1`

	getUserByEmail = `SELECT user_id, first_name, last_name, email, password, avatar, role, updated_at, created_at 
	FROM users WHERE email = $1`
)
