package types

import (
	"database/sql"
	"encoding/json"
)

// NullJSONString
type NullJSONString struct {
	sql.NullString
}

// MarshalJSON
func (n *NullJSONString) MarshalJSON() ([]byte, error) {
	if n.Valid {
		return json.Marshal(n.String)
	} else {
		return json.Marshal(nil)
	}
}

// UnmarshalJSON
func (n *NullJSONString) UnmarshalJSON(bytes []byte) error {
	var s *string
	if err := json.Unmarshal(bytes, &s); err != nil {
		return err
	}

	if s != nil {
		n.Valid = true
		n.String = *s
	} else {
		n.Valid = false
	}
	return nil
}
