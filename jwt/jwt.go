package jwt

import (
	"github.com/golang-jwt/jwt"
	"time"
)

// go get github.com/golang-jwt/jwt

type MyClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

var (
	mySigningKey = []byte("token_secert")
)

// 创建token
func CreateToken(username, password string) (string, error) {
	expireTime := time.Now().Add(time.Hour * 3)
	nowTime := time.Now()
	claims := MyClaims{
		Username: username,
		Password: password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			IssuedAt:  nowTime.Unix(),    // 创建时间
			Issuer:    "gin-web",         // 签发人
			Subject:   "user token",      // 主题
		},
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenString.SignedString(mySigningKey)
}

// 验证token
func CheckToken(tokenString string) (*MyClaims, bool) {
	token, _ := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	})
	if key, _ := token.Claims.(*MyClaims); !token.Valid {
		return key, true
	} else {
		return nil, false
	}
}
