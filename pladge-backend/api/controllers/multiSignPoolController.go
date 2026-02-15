package controllers

import (
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/request"
	"pledge-backend/api/models/response"
	"pledge-backend/api/services"
	"pledge-backend/api/validate"
	"pledge-backend/log"

	"github.com/gin-gonic/gin"
)

// MultiSignPoolController 多签池相关的接口控制器
type MultiSignPoolController struct {
}

// SetMultiSign 处理“设置多签配置”的请求
func (c *MultiSignPoolController) SetMultiSign(ctx *gin.Context) {
	res := response.Gin{Res: ctx}                     // 实例化响应工具类
	req := request.SetMultiSign{}                     // 定义接收请求参数的结构体
	log.Logger.Sugar().Info("SetMultiSign req ", req) // 记录入参日志（方便排错）

	// 1. 调用验证器：检查 ChainId 是否合法、必填项是否缺失
	errCode := validate.NewMutiSign().SetMultiSign(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil) // 如果校验不通过，直接返回错误给前端
		return
	}

	// 2. 调用业务层：将校验后的数据传给 Service 层进行数据库写入
	errCode, err := services.NewMutiSign().SetMultiSign(&req)
	if errCode != statecode.CommonSuccess {
		log.Logger.Error(err.Error())   // 业务执行出错（如数据库断连），记录详细 Error 日志
		res.Response(ctx, errCode, nil) // 返回服务器错误状态码
		return
	}

	// 3. 成功响应
	res.Response(ctx, statecode.CommonSuccess, nil)
	return
}

// GetMultiSign 处理“获取多签配置”的请求
func (c *MultiSignPoolController) GetMultiSign(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.GetMultiSign{}  // 定义获取配置所需的请求参数（如 ChainId）
	result := response.MultiSign{} // 定义返回给前端的数据结构
	log.Logger.Sugar().Info("GetMultiSign req ", nil)

	// 1. 调用验证器：校验 ChainId 等
	errCode := validate.NewMutiSign().GetMultiSign(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	// 2. 调用业务层：根据 ChainId 获取多签合约和账户信息
	// result 是引用传递，Service 会把查询到的数据填充到 result 里
	errCode, err := services.NewMutiSign().GetMultiSign(&result, req.ChainId)
	if errCode != statecode.CommonSuccess {
		log.Logger.Error(err.Error())
		res.Response(ctx, errCode, nil)
		return
	}

	// 3. 返回查询结果
	res.Response(ctx, statecode.CommonSuccess, result)
	return
}
