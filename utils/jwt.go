package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

const secretKey = "admin"

func GenerateToken(email string, userId int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  email,
		"userId": userId,
		"exp":    time.Now().Add(time.Hour * 2).Unix(),
	})
	return token.SignedString([]byte(secretKey))
}

func VerifyToken(token string) error {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		fmt.Println("Error was ", err.Error())
		return errors.New("Could not parse token.")
	}

	if !parsedToken.Valid {
		return errors.New("Invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("Invalid token claims")
	}
	email, ok := claims["email"].(string)
	userId, ok1 := claims["userId"]
	fmt.Println("userId was ", userId, reflect.TypeOf(userId))
	if !ok || !ok1 {
		fmt.Println("values types ", ok, ok1)
		return errors.New("Invalid token claims values")
	}
	fmt.Println("Email was ", email)
	fmt.Println("UserID was ", userId)
	return nil
}
