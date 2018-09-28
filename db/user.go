package db

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"fmt"
)

type User struct {
	ID    uuid.UUID `db:"id" json:"id,omitempty"`
	Email string    `db:"email" json:"email"`
	Hash  string    `db:"hash" json:"-"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c Credentials) Create() (User, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, err
	}

	u := User{ID: uuid.New(), Email: c.Email, Hash: string(h)}

	_, err = conn.Exec("INSERT INTO users (id, email, hash) VALUES (?, ?, ?)", u.ID, u.Email, u.Hash)

	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return u, nil
}
