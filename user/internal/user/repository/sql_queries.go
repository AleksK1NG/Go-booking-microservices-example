package repository

const (
	createUserQuery = `INSERT INTO users (first_name, last_name, email, password, avatar, role) 
	VALUES ($1,$2,$3,$4,$5,$6) 
	RETURNING user_id, first_name, last_name, email, avatar, role, updated_at, created_at`

	getUserByIDQuery = `SELECT user_id, first_name, last_name, email, avatar, role, updated_at, created_at FROM users WHERE user_id = $1`

	getUserByEmail = `SELECT user_id, first_name, last_name, email, password, avatar, role, updated_at, created_at 
	FROM users WHERE email = $1`

	updateUserQuery = `UPDATE users 
		SET first_name = COALESCE(NULLIF($1, ''), first_name), 
	    last_name = COALESCE(NULLIF($2, ''), last_name), 
	    email = COALESCE(NULLIF($3, ''), email), 
	    role = COALESCE(NULLIF($4, '')::role, role)
		WHERE user_id = $5
	    RETURNING user_id, first_name, last_name, email, role, avatar, updated_at, created_at`
)
