package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Models struct {
	User Users
}

type Users struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	Active    string    `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const dbOpsTimeout = 3 * time.Second

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{
		User: Users{},
	}
}

func (u *Users) Insert(user *Users) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbOpsTimeout)

	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		return "", err
	}

	newId := uuid.New().String()

	statement := `insert into users (id,email,password,name,active,created_at,updated_at) 
	values ($1,$2,$3,$4,$5,$6,$7) returning id`

	err = db.QueryRowContext(
		ctx,
		statement,
		newId,
		user.Email,
		hashedPassword,
		user.Name,
		true,
		time.Now(),
		time.Now(),
	).Scan(&newId)

	if err != nil {
		return "", nil
	}

	return newId, nil
}
