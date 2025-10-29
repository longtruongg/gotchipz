package cfg

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var count = 0 // 91  = task done

// SignaturePharosHub todo verify sign by etherscan it work, but server always return mismatch address :/
func SignaturePharosHub(ctx context.Context, provider *ethclient.Client, prik *ecdsa.PrivateKey) (interface{}, error) {
	ctxTime, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()
	nonce, err := provider.PendingNonceAt(ctxTime, ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("get nonce failed: %v", err)
	}
	signature, err := signSignature(provider, ctxTime, prik)
	if err != nil {
		return nil, fmt.Errorf("failed to sign signature: %w", err)
	}

	checksumAddr := common.HexToAddress(ADDRESS.Hex()).Hex()
	baseTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	payload := &PayloadSign{
		Address:   checksumAddr,
		ChainID:   "688689",
		Domain:    "testnet.pharosnetwork.xyz",
		Nonce:     fmt.Sprintf("%d", nonce),
		Signature: hexutil.Encode(signature),
		Timestamp: baseTime,
		Wallet:    "MetaMask",
	}
	fmt.Printf("Final Address: %s\n", payload.Address)
	fmt.Printf("Final Signature: %s\n", payload.Signature)
	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payloadSign %w", err)
	}
	pay := payHub{
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

func signSignature(provider *ethclient.Client, ctxTime context.Context, prik *ecdsa.PrivateKey) ([]byte, error) {
	baseTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	nonce, err := provider.PendingNonceAt(ctxTime, ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("get nonce failed: %v", err)
	}
	message := fmt.Sprintf(`testnet.pharosnetwork.xyz wants you to sign in with your Ethereum account:
%s

I accept the Pharos Terms of Service: testnet.pharosnetwork.xyz/privacy-policy/Pharos-PrivacyPolicy.pdf

URI: https://testnet.pharosnetwork.xyz

Version: 1

Chain ID: 688689

Nonce: %d

Issued At: %s`, ADDRESS.Hex(), nonce, baseTime)
	msgBytes := accounts.TextHash([]byte(message))
	signature, err := crypto.Sign(msgBytes, prik)
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %v", err)
	}
	// Ethereum expects the recovery byte to be 27 or 28
	if signature[64] < 27 {
		signature[64] += 27
	}
	signatureHex := hexutil.Encode(signature)
	if !strings.HasPrefix(signatureHex, "0x") {
		signatureHex = "0x" + signatureHex
	}
	return signature, nil
}
func SendNativePhrs(ctx context.Context, provider *ethclient.Client, des string, key *ecdsa.PrivateKey) (string, error) {
	ctxTime, cancle := context.WithTimeout(ctx, 20*time.Second)
	defer cancle()
	gasPrice, err := provider.SuggestGasPrice(ctxTime)
	if err != nil {
		return "", fmt.Errorf("get gas price failed: %v", err)
	}
	nonce, err := provider.PendingNonceAt(context.Background(), ADDRESS)
	if err != nil {
		return "", fmt.Errorf("get nonce failed: %v", err)
	}
	gasLimit := uint64(21000)
	defaults := big.NewInt(100000000000000)
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
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(688689)), key)
	if err != nil {
		return "", fmt.Errorf("sign tx failed: %v", err)
	}
	err = provider.SendTransaction(ctxTime, signTx)
	if err != nil {
		return "", fmt.Errorf("send tx failed: %v", err)
	}
	rep, err := bind.WaitMined(ctx, provider, signTx)
	if err != nil {
		return "", fmt.Errorf("check tx failed: %v", err)
	}
	if rep.Status == 0 {
		return "", fmt.Errorf("check tx failed: tx not confirmed %w", rep.TxHash)
	}
	return signTx.Hash().Hex(), nil
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
	pay := payHub{
		method: http.MethodPost,
		url:    PHAROS_CHECKIN,
		bear:   BEAR,
	}
	res, err := doRequest(context.Background(), pay, payloadJson)
	if err != nil {
		return fmt.Errorf("can not make request with payloadSign %w", err)
	}
	fmt.Println(string(res))
	return nil
}

func VerifyTransferTask(ctx context.Context, txHash string) (*PharosTaskResult, error) {
	ctxTime, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	pay := payHub{
		method: http.MethodPost,
		url:    VERIFY_TX_TRANSFER,
		bear:   BEAR,
	}
	phrs := PharosTask{
		Address: ADDRESS.String(),
		TaskId:  401,
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
	result := PharosTaskResult{}
	errj := json.Unmarshal(res, &result)
	if errj != nil {
		return nil, fmt.Errorf("cannot unmarshal result %w", err)
	}
	if result.Code == 0 {
		count++
	}
	return &result, nil
}
func verifySignatureBeforeSend(key *ecdsa.PrivateKey) bool {

	address := crypto.PubkeyToAddress(*key.Public().(*ecdsa.PublicKey))
	expectedAddress := common.HexToAddress("0xa8bf05d0881a225EB466175E37D60F96D65Ca1f6")

	fmt.Printf("Private key corresponds to address: %s\n", address.Hex())
	fmt.Printf("Expected address: %s\n", expectedAddress.Hex())

	return address.Hex() == expectedAddress.Hex()
}
