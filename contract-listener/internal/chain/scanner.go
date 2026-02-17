package chain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// 扫历史区块（生产必须）
func ScanHistory(
	client *ethclient.Client,
	addr common.Address,
	from, to uint64,
	out chan types.Log,
) error {

	q := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(from)),
		ToBlock:   big.NewInt(int64(to)),
		Addresses: []common.Address{addr},
	}

	logs, err := client.FilterLogs(context.Background(), q)
	if err != nil {
		return err
	}

	for _, l := range logs {
		out <- l
	}

	return nil
}
