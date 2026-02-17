package worker

import "github.com/ethereum/go-ethereum/core/types"

// worker池处理日志
func Start(n int, jobs chan types.Log, handler func(types.Log)) {

	for i := 0; i < n; i++ {

		go func() {

			for j := range jobs {
				handler(j)
			}

		}()
	}
}
