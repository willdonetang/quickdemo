package jwt

import (
	"fmt"
	"quickdemo/pkg/gredis"
	"quickdemo/pkg/setting"
	"strconv"
)

const (
	//后台token前缀
	UkuOfferBackendPre = "ukuoffer:backend:"
)

// 缓存token
func SetCacheToken(Userid int64, token string) error {
	//前缀id
	cacheKey := formatTokenKey(UkuOfferBackendPre, strconv.FormatInt(Userid, 10))
	err := gredis.Set(cacheKey, token)
	if err != nil {
		return err
	}
	return gredis.SetKeyExpire(cacheKey, setting.AppSetting.TokenExpireTime)
}

// 获取缓存的token
func GetCacheToken(useid int64) (string, error) {
	cacheKey := formatTokenKey(UkuOfferBackendPre, strconv.FormatInt(useid, 10))
	return gredis.GetStringValue(cacheKey)
}

// 检查用户token存在
func CheckToken(agency, username string) (bool, error) {
	cacheKey := formatTokenKey(agency, username)
	return gredis.CheckKey(cacheKey)
}

func formatTokenKey(pre, username string) string {
	return fmt.Sprintf("%s%s", "ukuoffer:backend:", username)
}
