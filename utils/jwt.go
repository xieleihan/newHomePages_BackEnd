package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"gin/config"
)

var JWT_SECTRET_KEY = []byte(config.SecretKey)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

/*
生成Token
*/
func GenerateToken(username string) (string, error) {
	expirationTime := time.Now().Add(2 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "gin-server",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECTRET_KEY)
}

/*
* 解析Token
*/
func ParseToken(tokenString string) (*Claims, error) {
	token,err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECTRET_KEY, nil
	})

	if claims,ok := token.Claims.(*Claims); ok && token.Valid {
		return claims,nil
	} 

	return nil, err
}