package db

import (
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type User struct {
	ID    int64  `json:"id,omitempty"`
	Email string `json:"email"`
	hash  string
}

// Convenience type for login use.
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c Credentials) Create() (User, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, err
	}

	res, err := conn.Exec(`INSERT INTO users (email, hash) VALUES (?, ?)`, c.Email, h)

	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	id, _ := res.LastInsertId()

	return User{ID: id, Email: c.Email}, nil
}
