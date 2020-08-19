package users

import (
	"crypto/rsa"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Claims struct {
	jwt.StandardClaims
	Roles map[uint64][]RoleType `json:"orgs"`
}

func HashPassword(pwd []byte, cost int) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword(pwd, cost)
}

func CheckPassword(hash, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
}

func ParseJwt(tokenString string, publicKey *rsa.PublicKey) (t *jwt.Token, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		log.WithError(err).Error("ParseWithClaims failed")
		return nil, err
	}
	if _, ok := token.Claims.(*Claims); ok && token.Valid {
		// Success!
		return token, nil
	}

	return nil, fmt.Errorf("failed to parse JWT claims")
}

func GetRequestJWTClaims(req *restful.Request) *Claims {
	claimAttr := req.Attribute("jwt.claims")
	claims, ok := claimAttr.(*Claims)
	if ok {
		return claims
	}
	return nil
}

func DecodeJWT(jwtRaw string, publicKey *rsa.PublicKey) *Claims {
	token, err := ParseJwt(jwtRaw, publicKey)
	if err != nil {
		return nil
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil
	}
	return claims
}

func CreateJWT(user *User, expirationDuration time.Duration) *Claims {
	jwtTime := time.Now()
	roles := make(map[uint64][]RoleType, len(user.Roles))
	for orgId, roleSet := range user.Roles {
		roles[orgId] = make([]RoleType, len(roleSet))
		for i := range roleSet {
			roles[orgId][i] = roleSet[i].Role
		}
	}
	return &Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwtTime.Add(expirationDuration).Unix(),
			IssuedAt:  jwtTime.Unix(),
			Subject:   user.Guid,
		},
		Roles: roles,
	}
}