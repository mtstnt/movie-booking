package helpers

import (
	"movie/gateway/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	SECRET = "mantap"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint32
}

func CreateUserToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		UserID: user.ID,
	})

	return token.SignedString([]byte(SECRET))
}

func CreateEmployeeToken(employee *models.Employee) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
		UserID: employee.ID,
	})

	return token.SignedString([]byte(SECRET))
}
