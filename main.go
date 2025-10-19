package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	BASE_URL         = "https://testnet.dplabs-internal.com"
	CONTRACT_ADDRESS = common.HexToAddress("0x0000000038f050528452d6da1e7aacfa7b3ec0a8")
	MINT_METHOD_ID   = "5b70ea9f" //avoid 0x
	ADDRESS          = common.HexToAddress("your-address")
)

func main() {
	provider, err := ethclient.Dial(BASE_URL)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	ctx := context.Background()
	prik, err := readKey()
	nonce, err := provider.PendingNonceAt(ctx, ADDRESS)
	if err != nil {
		log.Fatalf("Failed to get nonce: %v", err)
	}
	gasPrice, err := provider.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatalf("Failed to get gasprice: %v", err)
	}
	chainId, err := provider.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chainid: %v", err)
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &CONTRACT_ADDRESS,
		Data:     common.Hex2Bytes(MINT_METHOD_ID),
		Value:    big.NewInt(0),
		Gas:      300000,
		GasPrice: gasPrice,
	})
	signature, err := types.SignTx(tx, types.LatestSignerForChainID(chainId), prik)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
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
func readKey() (*ecdsa.PrivateKey, error) {
	//create .env file, prik=........
	str, err := os.ReadFile(".env")
	if err != nil {
		return nil, fmt.Errorf("cannot open file env")
	}
	return crypto.HexToECDSA(string(str)[5:])
}
