package repo

import (
	"contract-listener/internal/model"
	"sync"
)

var m sync.Map

// 幂等检查
func Exists(tx string, idx uint) bool {

	_, ok := m.Load(tx + string(rune(idx)))
	return ok
}

func Save(e model.Event) {
	m.Store(e.TxHash+string(rune(e.LogIndex)), true)
}
