package cfg

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var count = 0 // 91  = task done

// SignaturePharosHub todo verify sign by etherscan it work, but server always return mismatch address :/
func SignaturePharosHub(hub ParamHub) (interface{}, error) {
	ctxTime, cancel := context.WithTimeout(hub.Ctx, 20*time.Second)
	defer cancel()
	nonce, err := hub.Provider.PendingNonceAt(ctxTime, ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("get nonce failed: %v", err)
	}
	signature, err := signSignature(hub)
	if err != nil {
		return nil, fmt.Errorf("failed to sign signature: %w", err)
	}

	checksumAddr := common.HexToAddress(ADDRESS.Hex()).Hex()
	baseTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	payload := &PayloadSign{
		Address:   checksumAddr,
		ChainID:   "688689",
		Domain:    "testnet.pharosnetwork.xyz",
		Nonce:     strconv.FormatUint(nonce, 10),
		Signature: hexutil.Encode(signature),
		Timestamp: baseTime,
		Wallet:    "MetaMask",
	}
	log.Printf("Final Address: %s\n", payload.Address)
	log.Printf("Final Signature: %s\n", payload.Signature)
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payloadSign %w", err)
	}
	pay := PayHub{
		http.MethodPost,
		DOMAIN_PHAROS_SIGN_IN,
		"",
	}
	request, err := doRequest(ctxTime, pay, payloadJson)
	if err != nil {
		return nil, fmt.Errorf("cannot make request with payload %w", err)
	}
	var jwtHub ResponsePharosHub
	if err := json.Unmarshal(request, &jwtHub); err != nil {
		return nil, fmt.Errorf("cannot unmarshal response %w", err)
	}
	return jwtHub, nil
}

func signSignature(hub ParamHub) ([]byte, error) {

	baseTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	nonce, err := hub.Provider.PendingNonceAt(hub.Ctx, ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("get nonce failed: %v", err)
	}

	message := fmt.Sprintf(`testnet.pharosnetwork.xyz wants you to sign in with your Ethereum account:
%s

I accept the Pharos Terms of Service: testnet.pharosnetwork.xyz/privacy-policy/Pharos-PrivacyPolicy.pdf

URI: https://testnet.pharosnetwork.xyz

Version: 1

Chain ID: 688689

Nonce: %s

Issued At: %s`, ADDRESS.Hex(), strconv.FormatUint(nonce, 10), baseTime)
	msgBytes := accounts.TextHash([]byte(message))
	signature, err := crypto.Sign(msgBytes, hub.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %v", err)
	}
	// Ethereum expects the recovery byte to be 27 or 28
	if signature[64] < 27 {
		signature[64] += 27
	}
	return signature, nil
}
func SendNativePhrs(param ParamHub, des string) (string, error) {
	ctxTime, cancel := context.WithTimeout(param.Ctx, 20*time.Second)
	defer cancel()
	gasPrice, err := param.Provider.SuggestGasPrice(ctxTime)
	if err != nil {
		return "", fmt.Errorf("get gas price failed: %v", err)
	}
	nonce, err := param.Provider.PendingNonceAt(context.Background(), ADDRESS)
	if err != nil {
		return "", fmt.Errorf("get nonce failed: %v", err)
	}
	rand.NewSource(time.Now().UnixNano())
	gasLimit := uint64(21000)
	defaults := gachaPhrs()
	toDes := common.HexToAddress(des)
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			To:       &toDes,
			GasPrice: gasPrice,
			Value:    defaults,
			Gas:      gasLimit,
		},
	)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(688689)), param.Key)
	if err != nil {
		return "", fmt.Errorf("sign tx failed: %v", err)
	}
	err = param.Provider.SendTransaction(ctxTime, signTx)
	if err != nil {
		return "", fmt.Errorf("send tx failed: %v", err)
	}
	rep, err := bind.WaitMined(ctxTime, param.Provider, signTx)
	if err != nil {
		return "", fmt.Errorf("check tx failed: %v", err)
	}
	if rep.Status == 0 {
		return "", fmt.Errorf("check tx failed: tx not confirmed %w", rep.TxHash)
	}
	return signTx.Hash().Hex(), nil
}
func gachaPhrs() *big.Int {
	val := []int64{
		10000000000000,
		30000000000000,
		500000000000000,
		40000000000001,
		10000000000006,
		20000000000009,
		30000000020000,
	}
	rand.NewSource(time.Now().UnixNano())

	x := val[rand.Intn(len(val))]
	return big.NewInt(x)
}
func CheckInPharos(addr string) error {
	payload := struct {
		Address string `json:"address"`
	}{}
	payload.Address = addr
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("cannot marshal payloadSign %w", err)
	}
	pay := PayHub{
		method: http.MethodPost,
		url:    PHAROS_CHECKIN,
		bear:   BEAR,
	}
	_, err = doRequest(context.Background(), pay, payloadJson)
	if err != nil {
		return fmt.Errorf("can not make request with payloadSign %w", err)
	}
	return nil
}

func VerifyTransferTask(ctx context.Context, txHash string) (*PharosTaskResult, error) {

	ctxTime, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	pay := PayHub{
		method: http.MethodPost,
		url:    VERIFY_TX_TRANSFER,
		bear:   BEAR,
	}
	phrs := PharosTask{
		Address: ADDRESS.Hex(),
		TaskId:  401, //transfer task
		TxHash:  txHash,
	}
	payloadJson, err := json.Marshal(phrs)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payloadSign %w", err)
	}
	res, err := doRequest(ctxTime, pay, payloadJson)
	if err != nil {
		return nil, fmt.Errorf("can not make request with payloadSign %w", err)
	}

	var result PharosTaskResult
	err = json.Unmarshal(res, &result)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal result %w", err)
	}
	if result.Code == 0 {
		count++
	}

	return &result, nil
}
func ReadAddrs() []string {
	fil, err := os.Open("addr_dummy.txt")
	if err != nil {
		return nil
	}
	var str []string
	defer fil.Close()
	scanner := bufio.NewScanner(fil)
	for scanner.Scan() {
		str = append(str, scanner.Text())
	}
	return str
}
