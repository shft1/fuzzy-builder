package services

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type JWTIssuer interface {
	GenerateToken(userID int64, role string) (string, error)
	ParseToken(tokenString string) (userID int64, role string, err error)
}

type jwtIssuer struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTIssuer(secret, issuer string, ttl time.Duration) JWTIssuer {
	return &jwtIssuer{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

func (j *jwtIssuer) GenerateToken(userID int64, role string) (string, error) {
	claims := jwt.RegisteredClaims{
		ID:        strconv.FormatInt(userID, 10),
		Issuer:    j.issuer,
		Subject:   role,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

func (j *jwtIssuer) ParseToken(tokenString string) (int64, string, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := tok.Claims.(*jwt.RegisteredClaims)
	if !ok || !tok.Valid {
		return 0, "", errors.New("invalid token")
	}
	userID, _ := strconv.ParseInt(claims.ID, 10, 64)
	return userID, claims.Subject, nil
}

type PasswordHasher interface {
	HashPassword(plain string) (string, error)
	VerifyPassword(hash, plain string) error
}

type bcryptHasher struct{}

func NewPasswordHasher() PasswordHasher { return &bcryptHasher{} }

func (b *bcryptHasher) HashPassword(plain string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *bcryptHasher) VerifyPassword(hash, plain string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
