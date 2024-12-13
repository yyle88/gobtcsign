package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

type GetUtxoFromInterface interface {
	GetUtxoFrom(utxo wire.OutPoint) (*SenderAmountUtxo, error)
}

type SenderAmountUtxoClient struct {
	client *rpcclient.Client
}

func NewSenderAmountUtxoClient(client *rpcclient.Client) *SenderAmountUtxoClient {
	return &SenderAmountUtxoClient{client: client}
}

func (uc *SenderAmountUtxoClient) GetUtxoFrom(utxo wire.OutPoint) (*SenderAmountUtxo, error) {
	previousUtxoTx, err := GetRawTransaction(uc.client, utxo.Hash.String())
	if err != nil {
		return nil, errors.WithMessage(err, "get-raw-transaction")
	}
	previousOutput := previousUtxoTx.Vout[utxo.Index]

	previousAmount, err := btcutil.NewAmount(previousOutput.Value)
	if err != nil {
		return nil, errors.WithMessage(err, "get-previous-amount")
	}

	utxoFrom := NewSenderAmountUtxo(
		NewAddressTuple(previousOutput.ScriptPubKey.Address),
		int64(previousAmount),
	)
	return utxoFrom, nil
}

type SenderAmountUtxo struct {
	sender *AddressTuple
	amount int64
}

func NewSenderAmountUtxo(sender *AddressTuple, amount int64) *SenderAmountUtxo {
	return &SenderAmountUtxo{
		sender: sender,
		amount: amount,
	}
}

type SenderAmountUtxoCache struct {
	outputUtxoMap map[wire.OutPoint]*SenderAmountUtxo
}

func NewSenderAmountUtxoCache(utxoMap map[wire.OutPoint]*SenderAmountUtxo) *SenderAmountUtxoCache {
	return &SenderAmountUtxoCache{outputUtxoMap: utxoMap}
}

func (uc SenderAmountUtxoCache) GetUtxoFrom(utxo wire.OutPoint) (*SenderAmountUtxo, error) {
	utxoFrom, ok := uc.outputUtxoMap[utxo]
	if !ok {
		return nil, errors.Errorf("wrong utxo[%s:%d] not-exist-in-cache", utxo.Hash.String(), utxo.Index)
	}
	return utxoFrom, nil
}
