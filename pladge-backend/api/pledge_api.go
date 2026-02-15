package main

import (
	"pledge-backend/api/middlewares"
	"pledge-backend/api/models"
	"pledge-backend/api/models/kucoin"
	"pledge-backend/api/models/ws"
	"pledge-backend/api/routes"
	"pledge-backend/api/static"
	"pledge-backend/api/validate"
	"pledge-backend/config"
	"pledge-backend/db"

	"github.com/gin-gonic/gin"
)

// 负责协调数据库、验证器、长连接（WebSocket）、价格抓取器以及 Web 路由的启动。
func main() {
	//  1. 初始化存储层：建立 MySQL 和 Redis 的物理连接
	db.InitMysql()
	db.InitRedis()
	// 2. 自动同步数据库表结构：根据 models 定义的 struct 自动创建/更新表
	models.InitTable()

	// 3. 注册验证器：将 gin 与 go-playground-validator 绑定
	// 确保之前看到的 validate 层逻辑能正常拦截非法参数
	validate.BindingValidator()

	// 4. 启动异步服务（协程）
	// 启动 WebSocket 服务器，处理前端的实时推送需求
	go ws.StartServer()

	// 启动外部交易所价格抓取：从 Kucoin 实时获取 PLGR 价格并存入缓存
	go kucoin.GetExchangePrice()

	// 5. 启动 Gin Web 框架
	gin.SetMode(gin.ReleaseMode) // 生产模式：禁用冗余的调试日志
	app := gin.Default()         // 创建默认的 Gin 引擎

	// 6. 配置静态文件服务
	// 将服务器本地目录映射到 URL "/storage/"，以便前端能访问代币图标（Logo）
	staticPath := static.GetCurrentAbPathByCaller()
	app.Static("/storage/", staticPath)

	// 7. 注入中间件与路由
	app.Use(middlewares.Cors()) // 全局使用跨域中间件
	routes.InitRoute(app)       // 初始化所有 API 路由映射

	// 8. 启动监听：根据配置文件的端口号开启 HTTP 服务
	_ = app.Run(":" + config.Config.Env.Port)

}

/*
 If you change the version, you need to modify the following files'
 config/init.go
*/

// 模块名称,						职责描述,					核心逻辑
// 持久层 (DB/Models),存储中心,MySQL (结构化数据) + Redis (状态指纹/JWT/Session)。
// 同步层 (Schedule),链下索引,监听 BSC 链上状态，解决直接查链慢、查不到历史的问题。
// 接入层 (Middleware),门户安全,处理跨域、JWT 令牌解析、管理员权限校验。
// 逻辑层 (Service),业务大脑,处理多签配置、池信息聚合、Kucoin 价格换算。
// 交互层 (Web/WS),消息分发,RESTful API 供前端主动拉取，WebSocket 供服务器主动推送。
