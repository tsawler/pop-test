package models

import "time"

type User struct {
	ID int `db:"id"`
	FirstName string `db:"first_name"`
	LastName string `db:"last_name"`
	UserActive int `db:"user_active"`
	Email string `db:"email"`
	Password string `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
}
