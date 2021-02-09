package utils

import (
	"strconv"

	uuid "github.com/satori/go.uuid"
)

// CreateSQLPlaceholders Generate postgres $ placeholders
func CreateSQLPlaceholders(length int) string {
	var placeholders string
	for i := 0; i < length; i++ {
		placeholders += `$` + strconv.Itoa(i+1) + `,`

	}
	placeholders = placeholders[:len(placeholders)-1]
	return placeholders
}

// ConvertStringArrToUUID convert string slice to uuid
func ConvertStringArrToUUID(ids []string) ([]uuid.UUID, error) {
	uids := make([]uuid.UUID, 0, len(ids))
	for _, id := range ids {
		uid, err := uuid.FromString(id)
		if err != nil {
			return nil, err
		}
		uids = append(uids, uid)
	}

	return uids, nil
}
