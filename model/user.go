package model

import (
	"database/sql"
	"encoding/json"
)

// NullInt sql.NullInt represent nil or nil
type NullInt sql.NullInt64

// NullString sql.NullString represent nil or string
type NullString sql.NullString

// User struct in processing
type User struct {
	ID       NullInt    `json:"id"`
	Username NullString `json:"username"`
	Password NullString `json:"password"`
	Email    NullString `json:"email"`
}

// MarshalJSON custom json.Marshal method for NullInt
func (s NullInt) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return json.Marshal(s.Int64)
}

// Scan implementation of sql.Scanner
func (s *NullInt) Scan(value interface{}) error {
	var v sql.NullInt64
	if err := v.Scan(value); err != nil {
		return err
	}
	*s = NullInt(v)
	return nil
}

// UnmarshalJSON custom json.Unmarshal method for NullInt
func (s *NullInt) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.Int64)
	s.Valid = true
	if err != nil {
		s.Valid = false
	}
	return err
}

// MarshalJSON custom json.Marshal() for NullString
func (s NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return json.Marshal("")
	}
	return json.Marshal(s.String)
}

// Scan implement sql.Scanner for NullString
func (s *NullString) Scan(data interface{}) error {
	var str sql.NullString
	if err := str.Scan(data); err != nil {
		return err
	}
	*s = NullString(str)
	return nil
}

// UnmarshalJSON custom json.Marshal() for NullString
func (s *NullString) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &s.String)
	s.Valid = true
	if err != nil {
		s.Valid = false
	}
	return err
}
