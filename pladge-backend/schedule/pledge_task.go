package main

import (
	"pledge-backend/db"
	"pledge-backend/schedule/models"
	"pledge-backend/schedule/tasks"
)

func main() {

	// init mysql
	db.InitMysql()

	// init redis
	db.InitRedis()

	// create table
	models.InitTable()

	// pool task
	tasks.Task() // 主要负责定时监控链上任务

}

/*
 If you change the version, you need to modify the following files'
 config/init.go
*/

// 基础设施层: db (MySQL, Redis), log (Zap/Logrus).

// 模型层: models (GORM 映射, ABI 定义).

// 工具层: utils (高精度计算, 邮件发送, MD5, 环境加载).

// 业务层: services (同步池数据, 抓取价格, 监控余额, 更新图标).

// 调度层: tasks (基于 gocron 的任务分发).

// 部署层: systemd (服务自愈与后台运行).
