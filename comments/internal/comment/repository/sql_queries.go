package repository

const (
	createCommentQuery = `INSERT INTO comments (hotel_id, user_id, message, photos, rating) 
	VALUES ($1, $2, $3, $4, $5) RETURNING comment_id, hotel_id, user_id, message, photos, rating, created_at, updated_at`

	getCommByIDQuery = `SELECT comment_id, hotel_id, user_id, message, photos, rating, created_at, updated_at FROM comments WHERE comment_id = $1`

	updateCommentQuery = `UPDATE comments SET message = COALESCE(NULLIF($1, ''), message), rating = $2, photos = $3
	WHERE comment_id = $4
	RETURNING comment_id, hotel_id, user_id, message, photos, rating, created_at, updated_at`

	getTotalCountQuery = `SELECT count(comment_id) as total FROM comments WHERE hotel_id = $1`

	getCommentByHotelIDQuery = `SELECT comment_id, hotel_id, user_id, message, photos, rating, created_at, updated_at FROM comments
	WHERE hotel_id = $1 OFFSET $2 LIMIT $3`
)
