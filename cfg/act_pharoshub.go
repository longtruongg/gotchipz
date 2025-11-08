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
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

var subscribeABI = `[{"inputs":[{"internalType":"address","name":"uAddress","type":"address"},{"internalType":"uint256","name":"uAmount","type":"uint256"}],"name":"subscribe","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

func SignaturePharosHub(hub *ParamHub) (interface{}, error) {
	ctxTime, cancel := context.WithTimeout(hub.Ctx, 20*time.Second)
	defer cancel()
	defer hub.Provider.Close()
	nonce, err := hub.Provider.PendingNonceAt(ctxTime, ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("get nonce failed: %v", err)
	}
	baseTime := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	signature, err := dataSign(hub, baseTime, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to sign signature: %w", err)
	}

	payload := &PayloadSign{
		Address:   ADDRESS.Hex(),
		ChainID:   "688689",
		Domain:    "testnet.pharosnetwork.xyz",
		Nonce:     strconv.FormatUint(nonce, 10),
		Signature: hexutil.Encode(signature),
		Timestamp: baseTime,
		Wallet:    "MetaMask",
	}

	payloadJson, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal payloadSign %w", err)
	}
	pay := PayHub{
		http.MethodPost,
		DOMAIN_PHAROS_SIGN_IN,
		"",
	}
	response, err := doRequest(ctxTime, pay, payloadJson)
	if err != nil {
		return nil, fmt.Errorf("cannot make request with payload %w", err)
	}
	var jwtHub ResponsePharosHub
	err = response.Decode(&jwtHub)
	if err != nil {
		return nil, fmt.Errorf("cannot decode response with payload %w", err)
	}
	err = SaveBear(jwtHub.Data.Jwt)
	if err != nil {
		return nil, fmt.Errorf("cannot save jwt with payload %w", err)
	}
	return jwtHub.Msg, nil
}
func LoadBearToken() error {
	data, err := os.ReadFile(PHAROS_BEAR)
	if err != nil {
		return fmt.Errorf("cannot read file bearer token %w", err)
	}
	BEAR = string(data)
	log.Printf("Loaded existing token (length: %d)", len(BEAR))
	return nil
}
func SaveBear(data string) error {
	err := os.WriteFile(PHAROS_BEAR, []byte(data), 0600)
	if err != nil {
		return fmt.Errorf("cannot save bear_token: %w", err)
	}
	return nil
}
func dataSign(hub *ParamHub, baseTime string, nonce uint64) ([]byte, error) {
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

type dataHub struct {
	gasLimit, nonce uint64
	gasPrice        *big.Int
}

func fetchGas(param *ParamHub) (*dataHub, error) {
	gasPrice, err := param.Provider.SuggestGasPrice(param.Ctx)
	if err != nil {
		return nil, fmt.Errorf("get gas price failed: %v", err)
	}
	nonce, err := param.Provider.PendingNonceAt(param.Ctx, ADDRESS)
	if err != nil {
		return nil, fmt.Errorf("get nonce failed: %v", err)
	}

	defer param.Provider.Close()
	return &dataHub{
		gasLimit: uint64(66000),
		gasPrice: gasPrice,
		nonce:    nonce,
	}, nil

}
func SendNativePhrs(param *ParamHub, des string) (string, error) {
	hub, err := fetchGas(param)
	if err != nil {
		return "", fmt.Errorf("get gas price failed: %v", err)
	}
	defaults := gachaPhrs()
	toDes := common.HexToAddress(des)
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    hub.nonce,
			To:       &toDes,
			GasPrice: hub.gasPrice,
			Value:    defaults,
			Gas:      hub.gasLimit,
		},
	)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(688689)), param.Key)
	if err != nil {
		return "", fmt.Errorf("sign tx failed: %v", err)
	}
	<-time.After(time.Second * 10)
	log.Print("waiting 10s")
	err = param.Provider.SendTransaction(param.Ctx, signTx)
	if err != nil {
		return "", fmt.Errorf("send tx failed: %v", err)
	}
	rep, err := bind.WaitMined(param.Ctx, param.Provider, signTx)
	if err != nil {
		return "", fmt.Errorf("check tx failed: %v", err)
	}
	if rep.Status == 0 {
		return "", fmt.Errorf("check tx failed: tx not confirmed %w", rep.TxHash)
	}
	defer param.Provider.Close()
	return signTx.Hash().Hex(), nil
}

// todo
// got encoded calldata.
func AssetoSubcribe(param *ParamHub) error {

	hub, err := fetchGas(param)
	if err != nil {
		return fmt.Errorf("get gas value got : %v", err)
	}
	data := buildSubscribeCalldata(USDT_ATLTIC)
	toAddr := common.HexToAddress(CASH_ATLTIC)
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    hub.nonce,
			GasPrice: hub.gasPrice,
			Gas:      hub.gasLimit,
			To:       &toAddr,
			Value:    big.NewInt(0),
			Data:     data,
		})
	auth, err := bind.NewKeyedTransactorWithChainID(param.Key, big.NewInt(688689))
	if err != nil {
		return fmt.Errorf("get auth got : %v", err)
	}
	signTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return fmt.Errorf("authSign tx failed: %v", err)
	}
	err = param.Provider.SendTransaction(param.Ctx, signTx)
	if err != nil {
		return fmt.Errorf("send tx failed: %v", err)
	}
	mined, err := bind.WaitMined(param.Ctx, param.Provider, signTx)
	if err != nil {
		return fmt.Errorf("check tx failed: %v", err)
	}
	if mined.Status == 0 {
		return fmt.Errorf("check tx failed: tx not confirmed %w", mined.TxHash.Hex())
	}
	defer param.Provider.Close()
	return nil
}
func buildSubscribeCalldata(uAddressHex string) []byte {
	contractABI, err := abi.JSON(strings.NewReader(subscribeABI))
	if err != nil {
		log.Fatal(err)
	}
	uAmount := big.NewInt(10000000) // This equals 0x989680
	// Pack the function call
	data, err := contractABI.Pack("subscribe", common.HexToAddress(uAddressHex), uAmount)
	if err != nil {
		log.Fatal(err)
	}
	return data
}
func freshNonce(ctx context.Context, addr string) uint64 {
	return 0
}
func gachaPhrs() *big.Int {
	val := []int64{
		10000000000000,
		30000000000000,
		5000000000000,
		40000000000001,
		10000000000006,
		20000000000009,
		30000000020000,
		110000020000,
	}
	rand.NewSource(time.Now().UnixNano())
	return big.NewInt(val[rand.Intn(len(val))])
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
	checkTime := time.Now().UTC()
	err = SaveLastChecked(checkTime)
	if err != nil {
		return fmt.Errorf("save_last_checked failed: %v", err)
	}
	return nil
}
func SaveLastChecked(tim time.Time) error {
	current := tim.UTC().Format(time.RFC3339)
	return os.WriteFile(PHAROS_CHECKIN_TIME, []byte(current), 0644)
}

func ShouldCheckin() (bool, error) {
	txt, err := os.ReadFile(PHAROS_CHECKIN_TIME)
	if err != nil {
		return false, fmt.Errorf("read last checkin failed: %v", err)
	}
	lastTime, err := time.Parse(time.RFC3339, string(txt))
	if err != nil {
		return false, fmt.Errorf("parse last checkin failed: %v", err)
	}
	return time.Since(lastTime) >= 24*time.Hour, nil
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
	err = res.Decode(&result)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal result %w", err)
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
