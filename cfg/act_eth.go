package cfg

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)
func SignTxData(hub *ParamHub, des string, data []byte,value *big.Int) (*types.Transaction, error) {
  gasValue,err:=fetchGas(hub)
  if err!=nil {
    return nil, fmt.Errorf("cannot fetch gasHub %v", err)
  }

  toDes:=	common.HexToAddress(des)
  // arc chain
  gasLimit :=gasValue.gasPrice.Mul(gasValue.gasPrice, big.NewInt(21000))
  tx:=types.NewTx(
	&types.LegacyTx{
		Nonce:    gasValue.nonce,
		GasPrice: gasValue.gasPrice,
		Gas:      gasLimit.Uint64(),
		Value:    value,
		Data:     data,
		To:       &toDes,
	},
  )

  chainId,err:=hub.Provider.ChainID(context.Background())
  if err!=nil{
	return nil, fmt.Errorf("cannot get chainId")
  }
  signTx,err:= types.SignTx(tx, types.NewEIP155Signer(chainId), hub.Key)
  if err!=nil {
    return nil, fmt.Errorf("cannot sign tx")
  }
  err=hub.Provider.SendTransaction(hub.Ctx, signTx)
  if err!=nil {
    return nil, fmt.Errorf("cannot send signTx %w",err)
  }
  return signTx, nil
}
