package main

import (
	"context"
	"fmt"
	"log"
	"main/cfg"

	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	provider, err := ethclient.Dial(cfg.ATLANTIC)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	ctx := context.Background()
	prik, err := cfg.ReadKey()
	if err != nil {
		log.Fatalf("can not read prikey %s", err)
	}

	param := cfg.ParamHub{
		Ctx:      ctx,
		Provider: provider,
		Key:      prik,
	}

	for _, addr := range cfg.ReadAddrs() {

		time.AfterFunc(5*time.Second, func() {
			txHash, err := cfg.SendNativePhrs(param, addr)
			if err != nil {
				log.Fatalf("failed to send phrs %s", err)
			}
			task, err := cfg.VerifyTransferTask(ctx, txHash)
			if err != nil {
				log.Fatalf("failed to verify transfer %s", err)
			}
			fmt.Println(task)
		})
	}
	select {}
}
