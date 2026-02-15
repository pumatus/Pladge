package controllers

import (
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models"
	"pledge-backend/api/models/request"
	"pledge-backend/api/models/response"
	"pledge-backend/api/services"
	"pledge-backend/api/validate"
	"pledge-backend/config"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// PoolController 处理所有与质押池相关的 HTTP 请求
type PoolController struct {
}

// PoolBaseInfo 获取质押池的基础信息（如合约地址、固定利率等静态数据）
func (c *PoolController) PoolBaseInfo(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.PoolBaseInfo{}
	var result []models.PoolBaseInfoRes
	// 1. 参数校验
	errCode := validate.NewPoolBaseInfo().PoolBaseInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	// 2. 调用 Service 根据 ChainId 获取数据库中的基础信息
	errCode = services.NewPool().PoolBaseInfo(req.ChainId, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	// 3. 返回查询结果列表
	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

// PoolDataInfo 获取质押池的实时动态数据（如总借款额、剩余额度等）
func (c *PoolController) PoolDataInfo(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.PoolDataInfo{}
	var result []models.PoolDataInfoRes
	// 1. 参数校验
	errCode := validate.NewPoolDataInfo().PoolDataInfo(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	// 2. 调用 Service 获取实时数据
	errCode = services.NewPool().PoolDataInfo(req.ChainId, &result)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

// TokenList 获取协议支持的代币列表（遵循标准的 TokenList 格式）
func (c *PoolController) TokenList(ctx *gin.Context) {

	req := request.TokenList{}
	result := response.TokenList{}
	// 1. 校验 ChainId
	errCode := validate.NewTokenList().TokenList(ctx, &req)
	if errCode != statecode.CommonSuccess {
		ctx.JSON(200, map[string]string{
			"error": "chainId error",
		})
		return
	}
	// 2. 获取数据库中的代币信息
	errCode, data := services.NewTokenList().GetTokenList(&req)
	if errCode != statecode.CommonSuccess {
		ctx.JSON(200, map[string]string{
			"error": "chainId error",
		})
		return
	}
	// 3. 构建符合标准的 TokenList JSON 结构（包含版本号、Logo、时间戳等）
	var BaseUrl = c.GetBaseUrl()
	result.Name = "Pledge Token List"
	result.LogoURI = BaseUrl + "storage/img/Pledge-project-logo.png"
	result.Timestamp = time.Now()
	result.Version = response.Version{
		Major: 2,
		Minor: 16,
		Patch: 12,
	}
	for _, v := range data {
		result.Tokens = append(result.Tokens, response.Token{
			Name:     v.Symbol,
			Symbol:   v.Symbol,
			Decimals: v.Decimals,
			Address:  v.Token,
			ChainID:  v.ChainId,
			LogoURI:  v.Logo,
		})
	}

	ctx.JSON(200, result)
	return
}

// Search 根据关键字或条件搜索特定的质押池
func (c *PoolController) Search(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.Search{}
	result := response.Search{}
	// 1. 校验搜索参数（如分页、关键字）
	errCode := validate.NewSearch().Search(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}
	// 2. 执行搜索业务逻辑
	errCode, count, pools := services.NewSearch().Search(&req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	result.Rows = pools  // 搜索到的池子数据
	result.Count = count // 总条数（用于分页）
	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

// DebtTokenList 获取借款相关的代币列表
func (c *PoolController) DebtTokenList(ctx *gin.Context) {
	res := response.Gin{Res: ctx}
	req := request.TokenList{}

	errCode := validate.NewTokenList().TokenList(ctx, &req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	errCode, result := services.NewTokenList().DebtTokenList(&req)
	if errCode != statecode.CommonSuccess {
		res.Response(ctx, errCode, nil)
		return
	}

	res.Response(ctx, statecode.CommonSuccess, result)
	return
}

// GetBaseUrl 辅助函数：根据环境配置自动生成当前服务的基础访问路径（用于拼接 Logo 地址）
func (c *PoolController) GetBaseUrl() string {

	domainName := config.Config.Env.DomainName
	domainNameSlice := strings.Split(domainName, "")
	pattern := "\\d+" // 正则匹配数字
	// 如果域名开头是数字（如 IP 地址 127.0.0.1），则需要拼上端口号
	isNumber, _ := regexp.MatchString(pattern, domainNameSlice[0])
	if isNumber {
		return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + ":" + config.Config.Env.Port + "/"
	}
	// 如果是正式域名，则直接使用域名
	return config.Config.Env.Protocol + "://" + config.Config.Env.DomainName + "/"
}
