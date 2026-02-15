package models

// 以太坊智能合约 ABI (Application Binary Interface) 的结构体模型
type AbiJson struct {
	Status  string      `json:"status"`  //状态码
	Message string      `json:"message"` //提示信息
	Result  interface{} `json:"result"`  //核心内容数据
}

type Outputs struct {
	InternalType string `json:"internalType"` // Solidity 内部类型
	Name         string `json:"name"`
	Type         string `json:"type"`
}

type AbiData struct {
	Inputs          []interface{} `json:"inputs"` // 输入参数列表
	Payable         bool          `json:"payable,omitempty"`
	StateMutability string        `json:"stateMutability,omitempty"` // 状态可变性。例如 view（只读）、pure（不读不写）、nonpayable（不收币写入）、payable（收币写入）。
	Type            string        `json:"type"`                      // 类型。常见值有 function（函数）、constructor（构造函数）、event（事件）
	Anonymous       bool          `json:"anonymous,omitempty"`
	Name            string        `json:"name,omitempty"` // 函数或事件的名字（如 transfer 或 balanceOf）
	Constant        bool          `json:"constant,omitempty"`
	Outputs         []Outputs     `json:"outputs,omitempty"` // 函数执行后的返回结果列表
}
