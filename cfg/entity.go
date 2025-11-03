package cfg

import (
	"context"
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/ethclient"
)

type PharosTask struct {
	Address string `json:"address"`
	TaskId  int    `json:"task_id"` // 401 - transfer task
	TxHash  string `json:"tx_hash"`
}

type ParamHub struct {
	Ctx      context.Context
	Provider *ethclient.Client
	Key      *ecdsa.PrivateKey
}
type PharosTaskResult struct {
	Code int `json:"code"`
	Data struct {
		TaskId   int  `json:"task_id"`
		Verified bool `json:"verified"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type PayloadSign struct {
	Address   string `json:"address"`
	Signature string `json:"signature"`
	Wallet    string `json:"wallet"`
	Nonce     string `json:"nonce"`
	ChainID   string `json:"chain_id"`
	Timestamp string `json:"timestamp"`
	Domain    string `json:"domain"`
}
type ResponsePharosHub struct {
	Code int `json:"code"`
	Data struct {
		Jwt string `json:"jwt"`
	} `json:"data"`
	Msg string `json:"msg"`
}

type Gotchipx struct {
	Balance       string        `json:"balance"`
	Ids           []string      `json:"ids"`
	GotchipusInfo []interface{} `json:"gotchipusInfo"`
	TotalCount    int           `json:"totalCount"`
}
type PetOwner struct {
	Owner   string `json:"owner"`
	TokenId string `json:"tokenId"`
}
type PetOwnerResponse struct {
	Info struct {
		Name             string `json:"name"`
		Uri              string `json:"uri"`
		Story            string `json:"story"`
		Owner            string `json:"owner"`
		Collateral       string `json:"collateral"`
		CollateralAmount string `json:"collateralAmount"`
		Level            string `json:"level"`
		Status           int    `json:"status"`
		Evolution        int    `json:"evolution"`
		Locked           bool   `json:"locked"`
		Epoch            int    `json:"epoch"`
		Utc              int    `json:"utc"`
		Dna              struct {
			GeneSeed    string `json:"geneSeed"`
			RuleVersion int    `json:"ruleVersion"`
		} `json:"dna"`
		Bonding int    `json:"bonding"`
		Growth  int    `json:"growth"`
		Wisdom  int    `json:"wisdom"`
		Aether  int    `json:"aether"`
		Singer  string `json:"singer"`
		Nonces  string `json:"nonces"`
	} `json:"info"`
	TokenBoundAccount string `json:"tokenBoundAccount"`
	TokenName         string `json:"tokenName"`
}
