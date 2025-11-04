package main

import (
	"context"
	"fmt"
	"log"
	"main/cfg"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	provider, err := ethclient.Dial(cfg.PharosGotChips)
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
	res, err := cfg.GenSignatureGotchipus(param, cfg.CLAIM_WEARABLE)
	if err != nil {
		log.Fatalf("cannot getn 0 %s", err)
	}
	fmt.Println(res)
}
