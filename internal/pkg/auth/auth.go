package auth

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/klaital/volunteer-savvy-backend/internal/pkg/users"
	"golang.org/x/crypto/bcrypt"
	"time"
)
type Claims struct {
	jwt.StandardClaims
	Roles map[uint64][]users.Role `json:"orgs"`
}

func HashPassword(pwd []byte, cost int) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword(pwd, cost)
}

func CheckPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

func ParseJwt(tokenString string, publicKey *rsa.PublicKey) (t *jwt.Token, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(Claims); ok && token.Valid {
		// Success!
		return token, nil
	}

	return nil, fmt.Errorf("failed to parse JWT claims")
}

func CreateJWT(user *users.User, expirationDuration time.Duration) *Claims {
	jwtTime := time.Now()
	return &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtTime.Add(expirationDuration).Unix(),
			IssuedAt:  jwtTime.Unix(),
			Subject:   user.Guid,
		},
	}
}