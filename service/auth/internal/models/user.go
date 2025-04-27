package models

import "time"

type User struct {
	ID        int64     `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password_hash"`
	CreatedAt time.Time `db:"created_at"`
}
