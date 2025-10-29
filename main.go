package main

import (
	"context"
	"fmt"
	"log"
	"main/cfg"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	provider, err := ethclient.Dial("https://atlantic.dplabs-internal.com")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	ctx := context.Background()
	prik, err := cfg.ReadKey()
	if err != nil {
		log.Fatalf("can not read prikey %s", err)
	}
	hub, err := cfg.SignaturePharosHub(ctx, provider, prik)
	if err != nil {
		log.Fatalf("can not connect to Pharos Hub %v", err)
	}
	fmt.Println(hub)

}
