package main

import (
	"context"
	"log"

	"main/cfg"

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

	param := &cfg.ParamHub{
		Ctx:      ctx,
		Provider: provider,
		Key:      prik,
	}
	err = cfg.AssetoSubcribe(param)
	if err != nil {
		log.Fatalf("cannot getn 0 %s", err)
	}

}
