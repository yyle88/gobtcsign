package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

type GetUtxoFromInterface interface {
	GetUtxoFrom(utxo wire.OutPoint) (*UtxoSenderAmount, error)
}

type UtxoFromClient struct {
	client *rpcclient.Client
}

func NewUtxoFromClient(client *rpcclient.Client) *UtxoFromClient {
	return &UtxoFromClient{client: client}
}

func (uc *UtxoFromClient) GetUtxoFrom(utxo wire.OutPoint) (*UtxoSenderAmount, error) {
	preTxn, err := GetRawTransaction(uc.client, utxo.Hash.String())
	if err != nil {
		return nil, errors.WithMessage(err, "get-raw-txn")
	}
	preOut := preTxn.Vout[utxo.Index]

	preAmt, err := btcutil.NewAmount(preOut.Value)
	if err != nil {
		return nil, errors.WithMessage(err, "get-pre-amt")
	}

	utxoFrom := NewUtxoSenderAmount(
		NewAddressTuple(preOut.ScriptPubKey.Address),
		int64(preAmt),
	)
	return utxoFrom, nil
}

type UtxoSenderAmount struct {
	sender *AddressTuple
	amount int64
}

func NewUtxoSenderAmount(sender *AddressTuple, amount int64) *UtxoSenderAmount {
	return &UtxoSenderAmount{
		sender: sender,
		amount: amount,
	}
}

type OutPointUtxoSenderAmountMap struct {
	mxp map[wire.OutPoint]*UtxoSenderAmount
}

func NewOutPointUtxoSenderAmountMap(mxp map[wire.OutPoint]*UtxoSenderAmount) *OutPointUtxoSenderAmountMap {
	return &OutPointUtxoSenderAmountMap{mxp: mxp}
}

func (uc OutPointUtxoSenderAmountMap) GetUtxoFrom(utxo wire.OutPoint) (*UtxoSenderAmount, error) {
	utxoFrom, ok := uc.mxp[utxo]
	if !ok {
		return nil, errors.Errorf("not-exist-utxo[%s:%d]", utxo.Hash.String(), utxo.Index)
	}
	return utxoFrom, nil
}
