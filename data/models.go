package data

import (
	"context"
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	_ "crypto"
)

type Models struct {
	Users Users
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
var redisClient = &RedisClient{}
var jwtKey []byte
var jwtExpiration = time.Hour * 24 * 7 // 7 days

func New(dbPool *sql.DB, envJwtKey string) *Models {
	db = dbPool
	jwtKey = []byte(envJwtKey)
	return &Models{
		Users: Users{},
	}
}

func (_ *Users) Insert(user Users) (string, error) {
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
		return "", err
	}

	return newId, nil
}

func (_ *Users) FindByEmail(email string) (*Users, error) {
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

func (user *Users) GenerateJwt() (string, error) {
	randomRefreshToken := uuid.NewString()
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(jwtExpiration)
	claims["userId"] = user.ID
	claims["name"] = user.Name

	signedToken, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	err = redisClient.AddRefreshToken(user.ID, randomRefreshToken, jwtExpiration)

	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func (user *Users) VerifyJwt(token string) (string, bool) {
	parsedToken, err := jwt.NewParser().Parse(token, defaultKeyFunc)

	if err != nil {
		return "", false
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if !ok || !parsedToken.Valid {
		return "", false
	}

	refreshToken := claims["refreshToken"].(string)
	userId := claims["userId"].(string)

	redisClient.CheckRefreshTokenValid(userId, refreshToken)

	return refreshToken, true
}

func (user *Users) RefreshJwt(refreshToken string) (string, error) {
	token := jwt.New(jwt.SigningMethodES256)
	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(jwtExpiration)
	claims["userId"] = user.ID
	claims["name"] = user.Name

	signedToken, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	err = redisClient.AddRefreshToken(user.ID, refreshToken, jwtExpiration)

	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func defaultKeyFunc(t *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}
