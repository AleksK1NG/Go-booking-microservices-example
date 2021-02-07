package repository

const (
	getHotelByIDQuery = `SELECT hotel_id, email, name, location, description, comments_count, 
       	country, city, ((coordinates::POINT)[0])::decimal, ((coordinates::POINT)[1])::decimal, rating, photos, image, created_at, updated_at 
		FROM hotels WHERE hotel_id = $1`

	updateHotelQuery = `UPDATE hotels 
		SET email = COALESCE(NULLIF($1, ''), email), name = $2, location = $3, description = $4, 
	 	country = $5, city = $6, coordinates = ST_GeomFromEWKT($7)
		WHERE hotel_id = $8
	    RETURNING hotel_id, email, name, location, description, comments_count, 
       	country, city, ((coordinates::POINT)[0])::decimal, ((coordinates::POINT)[1])::decimal, rating, photos, image, created_at, updated_at`

	createHotelQuery = `INSERT INTO hotels (name, location, description, image, photos, coordinates, email, country, city, rating) 
	VALUES ($1, $2, $3, $4, $5, ST_GeomFromEWKT($6), $7, $8, $9, $10) RETURNING hotel_id, created_at, updated_at`

	getTotalHotelsCountQuery = `SELECT COUNT(*) as total FROM hotels`

	getHotelsQuery = `SELECT hotel_id, email, name, location, description, comments_count, 
       	country, city, ((coordinates::POINT)[0])::decimal, ((coordinates::POINT)[1])::decimal, rating, photos, image, created_at, updated_at 
       	FROM hotels OFFSET $1 LIMIT $2`
)
