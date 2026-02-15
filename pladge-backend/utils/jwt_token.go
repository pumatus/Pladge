package utils

import (
	"pledge-backend/config"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ✅ 时间处理 (time.go) -> 解决任务调度、日志记录。

// ✅ 基础转换 (functions.go) -> 解决 ID 生成、类型转换。

// ✅ 并发安全 (map.go) -> 解决内存缓存。

// ✅ 身份验证 (jwt_token.go) -> 解决用户权限。

// ✅ 通知系统 (email.go) -> 解决告警和通知

// 签发令牌
func CreateToken(username string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
	}) // 30天有效期
	//HS256 对令牌加密
	token, err := at.SignedString([]byte(config.Config.Jwt.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

// 解析令牌
func ParseToken(token string, secret string) (string, error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}) // 接收 验证
	if err != nil {
		return "", err
	}
	return claim.Claims.(jwt.MapClaims)["username"].(string), nil
}
