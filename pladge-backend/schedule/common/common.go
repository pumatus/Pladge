package common

import (
	"os"
	"pledge-backend/log"
)

var PlgrAdminPrivateKey string

// 从环境变量中读取敏感配置
// export plgr_admin_private_key="你的十六进制私钥"
func GetEnv() {

	var ok bool

	PlgrAdminPrivateKey, ok = os.LookupEnv("plgr_admin_private_key")
	if !ok {
		log.Logger.Error("environment variable is not set")
		panic("environment variable is not set")
	}

}
