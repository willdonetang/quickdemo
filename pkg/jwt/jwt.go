package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"quickdemo/pkg/logf"
	"quickdemo/pkg/setting"
	"time"
)

var jwtSecret []byte

func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

// 载荷，可以加一些自己需要的信息
type UkuOfferBackendClaims struct {
	ID          int64  `json:"userId"`
	AccountName string `json:"accountName"`
	jwt.StandardClaims
}

// GenerateToken generate tokens used for auth
func GenerateToken(claims UkuOfferBackendClaims) (string, error) {
	nowTime := time.Now()
	//expireTime := nowTime.Add(time.Duration(setting.AppSetting.TokenExpireTime) * time.Hour) // 1一个小时的过期时间
	//设置token过期时间
	claims.ExpiresAt = nowTime.Add(time.Duration(setting.AppSetting.TokenExpireTime)).Unix()
	claims.Issuer = "uku_backend"
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)

	err = SetCacheToken(claims.ID, token)
	if err != nil {
		logf.Error("token_cache.SetCacheToken:", err)
	}
	return token, err
}

// ParseToken parsing token
func ParseToken(token string) (*UkuOfferBackendClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &UkuOfferBackendClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*UkuOfferBackendClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
