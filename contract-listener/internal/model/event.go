package model

// 数据库存储的事件结构
type Event struct {
	TxHash   string
	Block    uint64
	LogIndex uint
}
