package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateToken(userID uint) (string, error) {
	nowTime := time.Now()
	// 设置为一个月 (30天)
	expireTime := nowTime.Add(30 * 24 * time.Hour)
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
			Issuer:    "media",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return "", err
	}
	fmt.Println("Token generated successfully")
	return tokenString, nil
}

func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.secret")), nil
	})
	if err != nil {
		return nil, err
	}
	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
