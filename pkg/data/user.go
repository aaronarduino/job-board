package data

import (
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        string    `db:"id"`
	Email     string    `db:"email"`
	Verified  bool      `db:"verified"`
	CreatedAt time.Time `db:"created_at"`
}

func NewUser(email string) *User {
	return &User{
		Email:    email,
		Verified: false,
	}
}

// Get User
func GetUserByEmail(email string, db *sqlx.DB) (User, error) {
	var user User

	err := db.Get(&user, "SELECT * FROM users WHERE email = $1", email)
	if err != nil {
		log.Println("did not find user")
		return user, err
	}
	log.Println("did find user")
	return user, nil
}

func GetUser(id string, db *sqlx.DB) (User, error) {
	var user User

	err := db.Get(&user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (user *User) ItemID() string {
	return user.ID
}

func (user *User) Hash(secret string) string {
	input := fmt.Sprintf(
		"%s:%s:%s:%s",
		user.ID,
		user.Email,
		user.CreatedAt.String(),
		secret,
	)

	hash := sha1.New()
	hash.Write([]byte(input))

	return string(base64.URLEncoding.EncodeToString(hash.Sum(nil)))
}

// Save
//
// If User does not exist in database create new user
func (u *User) SaveToDB(db *sqlx.DB) error {
	// If user does not exist in database create new user otherwise do nothing
	if user, err := GetUserByEmail(u.Email, db); err == nil {
		// If user exists return early
		log.Println("In SaveToDB user found...")
		*u = user
		return nil
	}

	query := `INSERT INTO users
    (email, verified)
    VALUES ($1, $2)
    RETURNING *`

	params := []interface{}{
		u.Email,
		false,
	}

	if err := db.QueryRowx(query, params...).StructScan(u); err != nil {
		return err
	}
	return nil
}

// Update User Verification
func (u *User) UpdateUserVerification(id string, verified bool, db *sqlx.DB) (sql.Result, error) {
	return db.Exec(
		"UPDATE users SET verified = $1 WHERE id = $2",
		verified, id,
	)

}
