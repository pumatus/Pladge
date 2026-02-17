package chain

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

// 自动重连client
func MustDial(url string) *ethclient.Client {

	for {
		c, err := ethclient.Dial(url)
		if err == nil {
			return c
		}
		log.Println("dial retry")
		time.Sleep(time.Second * 3)
	}
}
