package cfg

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var X_FILE = "x_ids.txt"

func GenSignature(ctx context.Context, provider *ethclient.Client, data string, prik *ecdsa.PrivateKey) (*types.Transaction, error) {
	ctxTime, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	nonce, err := provider.PendingNonceAt(ctxTime, ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}
	gasPrice, err := provider.SuggestGasPrice(ctxTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get gasprice: %v", err)
	}
	chainId, err := provider.ChainID(ctxTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get chainid: %v", err)
	}
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &CONTRACT_ADDRESS,
		Data:     common.Hex2Bytes(data),
		Value:    big.NewInt(0),
		Gas:      300000,
		GasPrice: gasPrice,
	})
	defer provider.Client().Close()
	return types.SignTx(tx, types.LatestSignerForChainID(chainId), prik)
}
func ReadKey() (*ecdsa.PrivateKey, error) {
	//create .env file, prik=........
	str, err := os.ReadFile(".env")
	if err != nil {
		return nil, fmt.Errorf("cannot open file env")
	}
	return crypto.HexToECDSA(string(str)[5:])
}

// get all gotchipus id
func FetchAllGotChipus(ctx context.Context, wallet string) (*Gotchipx, error) {
	gotchipx := &Gotchipx{}
	ownerUrl := fmt.Sprintf("https://gotchipus.com/api/tokens/gotchipus?owner=%s&includeGotchipusInfo=false", wallet)
	result, err := doRequest(ctx, "GET", ownerUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gotchipus: %v", err)
	}
	err = json.Unmarshal(result, &gotchipx)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal gotchipus: %v", err)
	}
	if err := saveIdsToFile(gotchipx.Ids); err != nil {
		return nil, fmt.Errorf("failed to save gotchipusId to file: %v", err)
	}
	return gotchipx, nil
}
func saveIdsToFile(ids []string) error {
	fil, err := os.OpenFile(X_FILE, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("failed to open x_ids.txt")
	}
	defer fil.Close()
	scanner := bufio.NewScanner(fil)
	for _, id := range ids {
		if scanner.Scan() {
			fisr := scanner.Text()
			if id == fisr {
				return fmt.Errorf("ids already exist %s", err)
			}
		}
		_, err := fmt.Fprintln(fil, id)
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}
	return nil
}
func doRequest(ctx context.Context, method, url string, data []byte) ([]byte, error) {
	client := &http.Client{}
	ctxTime, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctxTime, method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("can not make request")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %v", err)
	}
	defer res.Body.Close()
	return io.ReadAll(res.Body)
}
