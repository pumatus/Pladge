package service

import (
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type DepositLend struct {
	From       common.Address
	Token      common.Address
	Amount     *big.Int
	MintAmount *big.Int
}

func HandleDeposit(a abi.ABI, logg types.Log) {

	// topic[0] already matched

	var e DepositLend

	// indexed参数在topics里
	e.From = common.HexToAddress(logg.Topics[1].Hex())
	e.Token = common.HexToAddress(logg.Topics[2].Hex())

	// 非indexed在data
	err := a.UnpackIntoInterface(&e, "DepositLend", logg.Data)
	if err != nil {
		log.Println("decode fail", err)
		return
	}

	log.Println("Deposit event:",
		e.From.Hex(),
		e.Token.Hex(),
		e.Amount,
		e.MintAmount,
	)
}
