package auth

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	"golang.org/x/crypto/bcrypt"
)
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
	Roles map[uint64][]users.Role
}

func HashPassword(pwd []byte, cost int) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword(pwd, cost)
}

func CheckPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

func ParseJwt(tokenString, secret string) (userGuid string, t *jwt.Token, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return "", nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Success!
		if sub, ok := claims["sub"].(string); ok {
			return sub, token, nil
		}
		return "", token, fmt.Errorf("invalid sub: %v", claims["sub"])
	}

	return "", nil, fmt.Errorf("failed to parse JWT claims")
}
