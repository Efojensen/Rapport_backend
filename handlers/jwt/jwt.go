package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(emailOrUsername string) (string, error) {
	secretKey := os.Getenv("SECRET_KEY")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"emailOrUsername": emailOrUsername,
			"exp": time.Now().Add(time.Hour * 24).Unix(),
		},
	)
	tokenString, err := token.SignedString([]byte(secretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(token string) error {
	secretKey := os.Getenv("SECRET_KEY")
	ptrToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return err
	}

	if !ptrToken.Valid {
		return errors.New("invalid token")
	}

	return nil
}