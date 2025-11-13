package cfg

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func SignTxData(hub *ParamHub, des string, data []byte, value *big.Int) (*types.Transaction, error) {
	gasValue, err := fetchGas(hub)
	if err != nil {
		return nil, fmt.Errorf("cannot fetch gasHub %v", err)
	}

	toDes := common.HexToAddress(des)
	// arc chain
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    gasValue.nonce,
			GasPrice: gasValue.gasPrice,
			Gas:      gasValue.gasLimit,
			Value:    value,
			Data:     data,
			To:       &toDes,
		},
	)
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(gasValue.chainId), hub.Key)
	if err != nil {
		return nil, fmt.Errorf("cannot sign tx")
	}
	err = hub.Provider.SendTransaction(hub.Ctx, signTx)
	if err != nil {
		return nil, fmt.Errorf("cannot send signTx %w", err)
	}
	return signTx, nil
}
