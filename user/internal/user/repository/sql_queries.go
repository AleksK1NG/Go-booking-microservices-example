package repository

const (
	createUserQuery = `INSERT INTO users (first_name, last_name, email, password, avatar, role) 
	VALUES ($1,$2,$3,$4,$5,$6) 
	RETURNING user_id, first_name, last_name, email, avatar, role, updated_at, created_at`
)
