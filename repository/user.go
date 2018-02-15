package repository

import (
	"database/sql"
	"fmt"
	"encoding/json"
	"log"
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

// UserRepo interface for User Repository layer
type UserRepo interface {
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	GetByUsername(name string) (User, error)
	Update(user User) (User, error)
	Delete(user User) error
	Save(user User) (User, error)
}

// MarshalJSON custom json.Marshal method for NullInt
func (s *NullInt) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte{}, nil
	}
	return json.Marshal(s.Int64)
}

// Scan implementation of sql.Scanner
func (s *NullInt) Scan(value interface{}) error {
	var v sql.NullInt64
	if err := v.Scan(value); err != nil {
		log.Fatal(err)
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
func (s *NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return json.Marshal("")
	}
	return json.Marshal(s.String)
}

// Scan implement sql.Scanner for NullString
func (s *NullString) Scan(data interface{}) error {
	var str sql.NullString
	if err := str.Scan(data); err != nil {
		log.Fatal(err)
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

// UserRepoIml implement the UserRepo interface
type UserRepoIml struct {
	DB *sql.DB
}

func (u User) String() string {
	return fmt.Sprintln(u.ID, u.Username, u.Password, u.Email)
}

// GetAll implement UserRepo.GetAll()
func (repo *UserRepoIml) GetAll() ([]User, error) {
	rows, err := repo.DB.Query("SELECT id, username, password, email from user")
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()
	userList := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
			log.Fatal(err)
		}
		userList = append(userList, user)
	}
	return userList, nil
}

// GetByID implement UserRepo.GetByID
func (repo *UserRepoIml) GetByID(id int) (User, error) {
	row := repo.DB.QueryRow("SELECT id, username, password, email FROM user WHERE id=?", id)
	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		return User{}, err
	}
	return user, nil
}

// GetByUsername implement UserRepo.GetByUsername
func (repo *UserRepoIml) GetByUsername(name string) (User, error) {
	row := repo.DB.QueryRow("SELECT id, username, password, email FROM user WHERE username=?", name)
	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.Password); err != nil {
		return User{}, err
	}
	return user, nil
}

// Update implement UserRepo.Update
func (repo *UserRepoIml) Update(user User) (User, error) {
	_, err := repo.DB.Exec("UPDATE user set username=?,password=?,email=? WHERE id=?",
		user.Username, user.Password, user.Email, user.ID)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Delete implement UserRepo.Delete
func (repo *UserRepoIml) Delete(user User) error {
	_, err := repo.DB.Exec("DELETE FROM user WHERE id=?", user.ID)
	return err
}

// Save implement UserRepo.Save
func (repo *UserRepoIml) Save(user User) (User, error) {
	result, err := repo.DB.Exec("INSERT INTO user (username, password, email) VALUES (?, ?, ?)",
		user.Username, user.Password, user.Email)
	if err != nil {
		return User{}, err
	}

	id, insertErr := result.LastInsertId()
	if insertErr != nil {
		return User{}, insertErr
	}
	user.ID = NullInt{Int64: id, Valid: true}
	return user, nil
}