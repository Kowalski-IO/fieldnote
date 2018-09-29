package db

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"strings"
)

type User struct {
	ID    uuid.UUID `db:"id" json:"id,omitempty"`
	Username string    `db:"username" json:"username"`
	Hash  string    `db:"hash" json:"-"`
}

type Credentials struct {
	Username    string `json:"username"`
	Password string `json:"password"`
}

func (c Credentials) Create() (User, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(c.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, err
	}

	u := User{Username: strings.ToLower(c.Username), Hash: string(h)}

	_, err = conn.Exec("INSERT INTO users (username, hash) VALUES ($1, $2)", u.Username, u.Hash)

	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	return u, nil
}
