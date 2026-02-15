package controllers

import (
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/request"
	"pledge-backend/api/models/response"
	"pledge-backend/api/services"
	"pledge-backend/api/validate"
	"pledge-backend/db"

	"github.com/gin-gonic/gin"
)

// UserController 处理管理员身份验证相关的请求
type UserController struct {
}

// Login 处理用户登录请求
func (c *UserController) Login(ctx *gin.Context) {
	res := response.Gin{Res: ctx} // 初始化响应工具
	req := request.Login{}        // 接收用户名和密码的结构体
	result := response.Login{}    // 登录成功后返回给前端的数据（包含 Token）

	// 1. 参数校验：验证用户名密码是否为空
	errCode := validate.NewUser().Login(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	// 2. 业务逻辑处理：
	// Service 层会核对数据库密码，生成 JWT Token，并将登录状态存入 Redis
	errCode = services.NewUser().Login(&req, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	// 3. 登录成功，返回 Result（内含 Token 和用户信息）
	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

// Logout 处理用户登出请求
func (c *UserController) Logout(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	// 1. 从 Gin 上下文中获取用户名
	// 这个用户名是在 CheckToken 中间件解析 Token 后通过 ctx.Set("username", ...) 存入的
	usernameIntf, _ := ctx.Get("username")

	// 2. 核心注销逻辑：从 Redis 中删除该用户的登录状态键
	// 即使前端保留了 JWT Token，由于 Redis 中找不到对应的 key，中间件也会拦截后续请求
	_, _ = db.RedisDelete(usernameIntf.(string))

	// 3. 返回成功
	res.Response(ctx, statecode.CommonSuccess, nil)
	return
}
