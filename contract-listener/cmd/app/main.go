package main

import (
	"bytes"
	"context"
	"contract-listener/internal/chain"
	"contract-listener/internal/service"
	"contract-listener/internal/worker"
	"contract-listener/pkg/config"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func main() {

	// ============================
	// 1. 加载配置
	// ============================
	cfg := config.Load()

	// ============================
	// 2. 连接链节点（自动重连）
	// ============================
	client := chain.MustDial(cfg.RPCWS)

	// ============================
	// 3. 合约地址
	// ============================
	addr := common.HexToAddress(cfg.Contract)

	// ============================
	// 4. 加载 ABI
	// ============================
	abiBytes, err := os.ReadFile("abi.json")
	if err != nil {
		log.Fatal("read abi failed:", err)
	}

	parsedABI, err := abi.JSON(bytes.NewReader(abiBytes))
	if err != nil {
		log.Fatal("parse abi failed:", err)
	}

	// ============================
	// 5. 获取 DepositLend topic
	// ============================
	depositTopic := parsedABI.Events["DepositLend"].ID

	// ============================
	// 6. 日志 channel（必须大buffer）
	// ============================
	logs := make(chan types.Log, 10000)

	// ============================
	// 7. 启动 worker池
	// ============================
	worker.Start(cfg.WorkerNum, logs, func(v types.Log) {

		// 只处理 DepositLend
		if len(v.Topics) > 0 && v.Topics[0] == depositTopic {

			// 调用解析函数
			service.HandleDeposit(parsedABI, v)

		}
	})

	// ============================
	// 8. 历史补扫（生产必须）
	// ============================
	latest, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Println("scan history 0 ->", latest)

	err = chain.ScanHistory(client, addr, 0, latest, logs)
	if err != nil {
		log.Fatal(err)
	}

	// ============================
	// 9. 实时订阅（生产必须）
	// ============================
	go chain.SubscribeLoop(client, addr, logs)

	log.Println("listener started...")

	// 阻塞
	select {}
}
