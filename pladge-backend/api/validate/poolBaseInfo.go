package validate

import (
	"io"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/request"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PoolBaseInfo struct{}

func NewPoolBaseInfo() *PoolBaseInfo {
	return &PoolBaseInfo{}
}

// PoolBaseInfo 校验获取质押池基础信息的请求参数
func (v *PoolBaseInfo) PoolBaseInfo(c *gin.Context, req *request.PoolBaseInfo) int {
	// 1. 尝试将前端传来的 JSON/Query 数据绑定到结构体
	err := c.ShouldBind(req)
	// 2. 基础格式检查
	if err == io.EOF {
		return statecode.ParameterEmptyErr // 请求体不能为空
	} else if err != nil {
		// 如果绑定失败，解析具体的验证错误（比如必填项缺失）
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			if e.Field() == "ChainId" && e.Tag() == "required" {
				return statecode.ChainIdEmpty // 返回“链ID不能为空”的状态码
			}
		}
		return statecode.CommonErrServerErr
	}
	// 3. 核心业务校验：只允许 BSC 主网(56) 和 测试网(97)
	if req.ChainId != 97 && req.ChainId != 56 {
		return statecode.ChainIdErr // 返回“暂不支持该链”的状态码
	}

	return statecode.CommonSuccess // 校验通过
}
