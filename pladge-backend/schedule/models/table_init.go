package models

import "pledge-backend/db"

// 定时任务监听事件 保存db 持久化数据
func InitTable() {
	db.Mysql.AutoMigrate(&PoolBase{})
	db.Mysql.AutoMigrate(&PoolData{})
	db.Mysql.AutoMigrate(&RedisTokenInfo{})
	db.Mysql.AutoMigrate(&TokenInfo{})
}
