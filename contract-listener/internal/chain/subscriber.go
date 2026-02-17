package chain

import (
	"context"
	"log"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 实时订阅（自动重连）
func SubscribeLoop(
	client *ethclient.Client,
	addr common.Address,
	out chan types.Log,
) {

	for {

		ch := make(chan types.Log)

		q := ethereum.FilterQuery{
			Addresses: []common.Address{addr},
		}

		sub, err := client.SubscribeFilterLogs(context.Background(), q, ch)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		for {
			select {

			case err := <-sub.Err():
				log.Println("sub error", err)
				time.Sleep(time.Second)
				break

			case v := <-ch:
				out <- v
			}
		}
	}
}
