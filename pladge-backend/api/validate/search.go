package validate

import (
	"io"
	"pledge-backend/api/common/statecode"
	"pledge-backend/api/models/request"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Search 定义搜索验证结构体
type Search struct{}

// NewSearch 初始化搜索验证器
func NewSearch() *Search {
	return &Search{}
}

// Search 执行具体的搜索参数校验逻辑
func (s *Search) Search(c *gin.Context, req *request.Search) int {
	// 1. 尝试将请求体中的 JSON 数据绑定到 req 结构体
	// ShouldBindJSON 会自动解析 Content-Type 为 application/json 的数据
	err := c.ShouldBindJSON(req)
	// 2. 检查请求体是否完全为空
	if err == io.EOF {
		return statecode.ParameterEmptyErr
	} else if err != nil {
		// 3. 字段校验：如果 JSON 解析成功但数据不符合 struct tag 定义（如 required）
		// 则进行类型断言，遍历具体的验证错误
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			// 如果是 ChainID 字段缺失
			if e.Field() == "ChainID" && e.Tag() == "required" {
				return statecode.ChainIdEmpty
			}
		}
		// 其他绑定错误（如 JSON 格式非法、类型不匹配等）返回服务器通用错误
		return statecode.CommonErrServerErr
	}
	// 4. 业务边界校验：只允许在 97 (BSC测试网) 或 56 (BSC主网) 进行搜索
	if req.ChainID != 97 && req.ChainID != 56 {
		return statecode.ChainIdErr
	}
	// 5. 校验通过，返回成功代码
	return statecode.CommonSuccess
}
