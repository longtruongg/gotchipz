package main

import (
	"context"
	"fmt"
	"log"
	"main/cfg"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	provider, err := ethclient.Dial(cfg.BaseUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	ctx := context.Background()
	prik, err := cfg.ReadKey()
	if err != nil {
		log.Fatalf("can not read prikey %s", err)
	}
	signature, err := cfg.GenSignature(ctx, provider, cfg.PET, prik)
	if err != nil {
		log.Fatalf("can not sign signature %s", err)
	}
	err = provider.SendTransaction(ctx, signature)
	if err != nil {
		log.Fatalf("cannot send tx %v", err)
	}
	rep, err := bind.WaitMined(ctx, provider, signature)
	if err != nil {
		log.Fatalf("failed to wait tx %v", err)
	}
	if rep.Status == 0 {
		fmt.Println("Tx FAILED (reverted). Check logs or use Tenderly for decode.")
	} else {
		fmt.Println("Tx SUCCESS! Minted. ", signature.Hash().String())
	}
}
