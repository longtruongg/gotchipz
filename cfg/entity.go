package cfg

import (
	"fmt"
	"time"
)

// data to signature
func dataSign(wallet, nonce, issued string) string {
	return fmt.Sprintf("testnet.pharosnetwork.xyz wants you to sign in with your Ethereum account:\n%s\n\nI accept the Pharos Terms of Service: testnet.pharosnetwork.xyz/privacy-policy/Pharos-PrivacyPolicy.pdf\n\nURI: https://testnet.pharosnetwork.xyz\n\nVersion: 1\n\nChain ID: 688688\n\nNonce: %s\n\nIssued At: %s\nfor  https://testnet.pharosnetwork.xyz",
		wallet, nonce, issued)
}

type PayloadSign struct {
	Address    string    `json:"address"`
	ChainId    string    `json:"chain_id"`
	Domain     string    `json:"domain"`      // domain pharos
	InviteCode string    `json:"invite_code"` // hardcode
	Nonce      string    `json:"nonce"`
	Signature  string    `json:"signature"` // gen signature
	Timestamp  time.Time `json:"timestamp"`
	Wallet     string    `json:"wallet" ` //defaultOkx
}
type Response struct {
	Code int `json:"code"`
	Data struct {
		Jwt string `json:"jwt"`
	} `json:"data"`
	Msg string `json:"msg"`
}
type SignIn struct {
	Address string `json:"address"`
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
