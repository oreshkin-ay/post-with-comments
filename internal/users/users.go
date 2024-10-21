package users

import (
	"database/sql"
	"errors"

	database "github.com/oreshkin/posts/internal/pkg/db/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"name"`
	Password string `json:"password"`
}

// Create creates a new user in the database
func (user *User) Create() error {
	statement, err := database.Db.Prepare("INSERT INTO users(username, password) VALUES($1, $2)")
	if err != nil {
		return err
	}
	defer statement.Close()

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return err
	}

	_, err = statement.Exec(user.Username, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

// GetUserIdByUsername returns the ID of a user by their username
func GetUserIdByUsername(username string) (int, error) {
	statement, err := database.Db.Prepare("SELECT id FROM users WHERE username = $1")
	if err != nil {
		return 0, err
	}
	defer statement.Close()

	var id int
	row := statement.QueryRow(username)
	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	return id, nil
}

// HashPassword hashes the given password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash checks if the provided password matches the hashed password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Authenticate authenticates the user by verifying the username and password
func (user *User) Authenticate() (bool, error) {
	statement, err := database.Db.Prepare("SELECT password FROM users WHERE username = $1")
	if err != nil {
		return false, err
	}
	defer statement.Close()

	var hashedPassword string
	row := statement.QueryRow(user.Username)
	err = row.Scan(&hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	isValid := CheckPasswordHash(user.Password, hashedPassword)
	return isValid, nil
}
