package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Models struct {
	Users *Users
}

type Users struct {
	ID                     string    `json:"id"`
	Email                  string    `json:"email"`
	Name                   string    `json:"name"`
	Password               string    `json:"-"`
	AuthenticatorSecretKey string    `json:"-"`
	Is2FAEnabled           bool      `json:"is2FAEnabled"`
	Role                   string    `json:"role"`
	Active                 bool      `json:"active"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}

type JwtUsers struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

const dbOpsTimeout = 3 * time.Second

var db *sql.DB

func New() *Models {
	if db == nil {
		db = connectDB()
		initAdminAccount()
	}

	return &Models{
		Users: &Users{},
	}
}

func (user *Users) ToJwtUser() *JwtUsers {
	return &JwtUsers{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}
}

func (user *JwtUsers) GetClaims() *map[string]any {
	return &map[string]any{
		"userID": user.ID,
		"email":  user.Email,
		"name":   user.Name,
		"role":   user.Role,
	}
}

func (*Users) ParseFromClaims(claims *jwt.MapClaims) *JwtUsers {
	return &JwtUsers{
		ID:    fmt.Sprintf("%s", (*claims)["userID"]),
		Email: fmt.Sprintf("%s", (*claims)["email"]),
		Name:  fmt.Sprintf("%s", (*claims)["name"]),
		Role:  fmt.Sprintf("%s", (*claims)["role"]),
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
	user.Is2FAEnabled = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	statement := `insert into users 
	(id,email,password,name,authenticatorsecretkey,is2faenabled,role,active,created_at,updated_at) 
	values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) returning id`

	err = db.QueryRowContext(
		ctx,
		statement,
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.AuthenticatorSecretKey,
		&user.Is2FAEnabled,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	).Scan(
		&newId,
	)

	if err != nil {
		return errors.New("unexpected error")
	}

	return nil
}

func (*Users) FindByEmail(email string) (*Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbOpsTimeout)

	defer cancel()

	query := `select 
	id,email,password,name,authenticatorsecretkey,is2faenabled,role,active,created_at,updated_at
	from users where email = $1`

	var user Users
	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.AuthenticatorSecretKey,
		&user.Is2FAEnabled,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (*Users) FindByUserID(id string) (*Users, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbOpsTimeout)

	defer cancel()

	query := `select 
	id,email,password,name,authenticatorsecretkey,is2faenabled,role,active,created_at,updated_at
	from users where id = $1`

	var user Users
	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.AuthenticatorSecretKey,
		&user.Is2FAEnabled,
		&user.Role,
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
