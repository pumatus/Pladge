package chain

import (
	"bytes"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// 加载ABI文件
func LoadABI(path string) abi.ABI {

	b, _ := os.ReadFile(path)

	parsed, _ := abi.JSON(
		bytes.NewReader(b),
	)

	return parsed
}

// 获取事件索引
func EventID(a abi.ABI, name string) [32]byte {
	return a.Events[name].ID
}
