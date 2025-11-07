package cfg

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var X_FILE = "x_ids.txt"

// todo change logic after
func GenSignatureGotchipus(param *ParamHub, data string) (bool, error) {
	ctxTime, cancel := context.WithTimeout(param.Ctx, 5*time.Second)
	defer cancel()
	nonce, err := param.Provider.PendingNonceAt(ctxTime, ADDRESS)
	if err != nil {
		return false, fmt.Errorf("failed to get nonce: %v", err)
	}
	gasPrice, err := param.Provider.SuggestGasPrice(ctxTime)
	if err != nil {
		return false, fmt.Errorf("failed to get gasprice: %v", err)
	}
	chainId, err := param.Provider.ChainID(ctxTime)
	if err != nil {
		return false, fmt.Errorf("failed to get chainid: %v", err)
	}
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &CONTRACT_ADDRESS,
		Data:     common.Hex2Bytes(data),
		Value:    big.NewInt(0),
		Gas:      300000,
		GasPrice: gasPrice,
	})
	signTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainId), param.Key)
	if err != nil {
		return false, fmt.Errorf("failed to sign tx: %v", err)
	}
	err = param.Provider.SendTransaction(ctxTime, signTx)
	if err != nil {
		return false, fmt.Errorf("failed to send tx: %v", err)
	}
	mined, err := bind.WaitMined(ctxTime, param.Provider, signTx.Hash())
	if err != nil {
		return false, fmt.Errorf("failed to mined tx: %v", err)
	}
	return mined.Status == 1, nil
}
func ReadKey() (*ecdsa.PrivateKey, error) {
	str, err := os.ReadFile(".env")
	if err != nil {
		return nil, fmt.Errorf("cannot open file env")
	}
	return crypto.HexToECDSA(string(str))
}

// get all gotchipus ids
func FetchAllGotChipus(ctx context.Context, wallet string) (*Gotchipx, error) {
	gotchipx := &Gotchipx{}

	ownerUrl := fmt.Sprintf("https://gotchipus.com/api/tokens/gotchipus?owner=%s&includeGotchipusInfo=false", wallet)
	pay := PayHub{
		http.MethodGet,
		ownerUrl,
		"",
	}

	result, err := doRequest(ctx, pay, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gotchipus: %v", err)
	}
	err = result.Decode(&gotchipx)
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

type PayHub struct {
	method string
	url    string
	bear   string
}

func doRequest(ctx context.Context, pay PayHub, data []byte) (*json.Decoder, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest(pay.method, pay.url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("can not make request %w", err)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	log.Printf("Raw JSON being sent: %s\n", string(data))
	if pay.bear != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", pay.bear))
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request failed: %v", err)
	}
	decoder := json.NewDecoder(res.Body)
	if !decoder.More() {
		return nil, fmt.Errorf("no JSON data in response %w", res.StatusCode)
	}
	defer res.Body.Close()
	return decoder, nil
}

var header = map[string]string{
	"Accept":         "application/json,text/plain, */*",
	"Origin":         "https://testnet.pharosnetwork.xyz",
	"referer":        "https://testnet.pharosnetwork.xyz/",
	"Sec-Fetch-Dest": "empty",
	"Sec-Fetch-Mode": "cors",
	"Sec-Fetch-Site": "same-origin",
	"Cache-Control":  "no-cache",
	"Content-type":   "application/json; charset=utf-8",
	"User-Agent":     "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Mobile Safari/537.36",
}
