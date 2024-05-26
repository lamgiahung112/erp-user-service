package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Models struct {
	Users *Users
}

type Users struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Password  string    `json:"-"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const dbOpsTimeout = 3 * time.Second

var db *sql.DB

func New(dbPool *sql.DB) *Models {
	db = dbPool

	return &Models{
		Users: &Users{},
	}
}

func (user *Users) GetClaims() *map[string]any {
	return &map[string]any{
		"userID": user.ID,
		"email":  user.Email,
		"name":   user.Name,
	}
}

func (*Users) Insert(user *Users) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbOpsTimeout)

	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)

	if err != nil {
		return err
	}

	newId := uuid.New().String()

	user.ID = newId
	user.Active = true
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	statement := `insert into users (id,email,password,name,active,created_at,updated_at) 
	values ($1,$2,$3,$4,$5,$6,$7) returning id`

	err = db.QueryRowContext(
		ctx,
		statement,
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	).Scan(
		&newId,
	)

	if err != nil {
		return err
	}

	return nil
}

func (*Users) FindByEmail(email string) (*Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbOpsTimeout)

	defer cancel()

	query := `select id, email, name, password, active, created_at, updated_at from users where email = $1`

	var user Users
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (user *Users) PasswordMatches(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(plainPassword))
	return err == nil
}
