package repository

import (
	"database/sql"
	"github.com/aceforeverd/enigma/model"
)

type User model.User

// UserRepo interface for User Repository layer
type UserRepo interface {
	InitTable() error
	GetAll() ([]User, error)
	GetByID(id int) (User, error)
	GetByUsername(name string) (User, error)
	Update(user User) (User, error)
	Delete(user User) error
	Save(user User) (User, error)
}

// UserRepoIml implement the UserRepo interface
type UserRepoIml struct {
	DB *sql.DB
}

// InitTable create user table if not exists
func (repo *UserRepoIml) InitTable() error {
	stm := `CREATE TABLE IF NOT EXISTS user (
			id INT NOT NULL AUTO_INCREMENT,
			username VARCHAR(255) NOT NULL,
			password VARCHAR(255) NOT NULL,
			email VARCHAR(255),
			PRIMARY KEY (id)) ENGINE=InnoDB`
	if _, err := repo.DB.Exec(stm); err != nil {
		return err
	}
	return nil
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
			return []User{}, err
		}
		userList = append(userList, user)
	}
	return userList, nil
}

// GetByID implement UserRepo.GetByID
func (repo *UserRepoIml) GetByID(id int) (User, error) {
	row := repo.DB.QueryRow("SELECT id, username, password, email FROM user WHERE id=?", id)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
		return User{}, err
	}
	return user, nil
}

// GetByUsername implement UserRepo.GetByUsername
func (repo *UserRepoIml) GetByUsername(name string) (User, error) {
	row := repo.DB.QueryRow("SELECT id, username, password, email FROM user WHERE username=?", name)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
		return User{}, err
	}
	return user, nil
}

// Update implement UserRepo.Update
func (repo *UserRepoIml) Update(user User) (User, error) {
	_, err := repo.DB.Exec("UPDATE user set username=?,password=?,email=? WHERE id=?",
		user.Username.String, user.Password.String, user.Email.String, user.ID.Int64)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

// Delete implement UserRepo.Delete
func (repo *UserRepoIml) Delete(user User) error {
	_, err := repo.DB.Exec("DELETE FROM user WHERE id=?", user.ID.Int64)
	return err
}

// Save implement UserRepo.Save
func (repo *UserRepoIml) Save(user User) (User, error) {
	result, err := repo.DB.Exec("INSERT INTO user (username, password, email) VALUES (?, ?, ?)",
		user.Username.String, user.Password.String, user.Email.String)
	if err != nil {
		return User{}, err
	}

	id, insertErr := result.LastInsertId()
	if insertErr != nil {
		return User{}, insertErr
	}
	user.ID = model.NullInt{Int64: id, Valid: true}
	return user, nil
}
