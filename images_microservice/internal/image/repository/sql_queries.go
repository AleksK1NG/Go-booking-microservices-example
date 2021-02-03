package repository

const (
	createImageQuery = `INSERT INTO images (image_url, is_uploaded) 
	VALUES ($1, $2) 
	RETURNING image_id, image_url, is_uploaded, created_at`
)
