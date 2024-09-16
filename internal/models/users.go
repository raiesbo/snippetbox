package models

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Define a new User struct. Notice how the field names and types align
// with the columns in the database "users" table?
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Definne a enw UserModel structu which wraps a database connection pool.
type UserModel struct {
	DB *sql.DB
}

// We'll user the Insert method to adda  new record to the "users" table
func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
VALUES ($1, $2, $3, NOW()) RETURNING id;`

	id := 0
	err = m.DB.QueryRow(stmt, name, email, string(hashedPassword)).Scan(&id)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" {
				return ErrDuplicateEamil
			}
		}
		return err
	}

	return nil
}

// We'll use the Authenticate method to veridfy whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// We'll use the Exists method to check if a user exists with a specific ID.
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}
