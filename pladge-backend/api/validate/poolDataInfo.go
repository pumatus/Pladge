package validate

import (
	"io"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/request"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// PoolDataInfo 定义质押池动态数据验证结构体
type PoolDataInfo struct{}

// NewPoolDataInfo 构造函数，用于初始化验证对象
func NewPoolDataInfo() *PoolDataInfo {
	return &PoolDataInfo{}
}

// PoolDataInfo 执行具体的参数验证逻辑
// c: Gin上下文，用于获取请求数据
// req: 绑定后的请求结构体指针
func (v *PoolDataInfo) PoolDataInfo(c *gin.Context, req *request.PoolDataInfo) int {
	// 1. 参数绑定：尝试将 HTTP 请求中的参数（Query 或 JSON）填充到 req 结构体中
	err := c.ShouldBind(req)
	// 2. 检查请求体是否为空 (End Of File)
	if err == io.EOF {
		return statecode.ParameterEmptyErr // 返回“参数为空”的状态码
	} else if err != nil {
		// 3. 字段级别验证：如果绑定失败，解析具体的 validator 错误
		// 这里使用了 go-playground/validator 库提供的错误断言
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			// 如果是 ChainId 字段且未通过 required（必填）标签校验
			if e.Field() == "ChainId" && e.Tag() == "required" {
				return statecode.ChainIdEmpty // 返回“链ID缺失”状态码
			}
		}
		// 其他未知绑定错误统一返回服务器异常
		return statecode.CommonErrServerErr
	}
	// 4. 业务逻辑校验：限制只能查询特定的区块链网络
	// 97: BSC Testnet (测试网)
	// 56: BSC Mainnet (主网)
	if req.ChainId != 97 && req.ChainId != 56 {
		return statecode.ChainIdErr
	}
	// 5. 校验全部通过
	return statecode.CommonSuccess
}
