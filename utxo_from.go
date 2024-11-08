package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

type GetUtxoFromInterface interface {
	GetUtxoFrom(utxo wire.OutPoint) (*UtxoSenderAmountTuple, error)
}

type UtxoFromClient struct {
	client *rpcclient.Client
}

func NewUtxoFromClient(client *rpcclient.Client) *UtxoFromClient {
	return &UtxoFromClient{client: client}
}

func (uc *UtxoFromClient) GetUtxoFrom(utxo wire.OutPoint) (*UtxoSenderAmountTuple, error) {
	preTxn, err := GetRawTransaction(uc.client, utxo.Hash.String())
	if err != nil {
		return nil, errors.WithMessage(err, "get-raw-txn")
	}
	preOut := preTxn.Vout[utxo.Index]

	preAmt, err := btcutil.NewAmount(preOut.Value)
	if err != nil {
		return nil, errors.WithMessage(err, "get-pre-amt")
	}

	utxoFrom := NewUtxoSenderAmountTuple(
		NewAddressTuple(preOut.ScriptPubKey.Address),
		int64(preAmt),
	)
	return utxoFrom, nil
}

type UtxoSenderAmountTuple struct {
	sender *AddressTuple
	amount int64
}

func NewUtxoSenderAmountTuple(sender *AddressTuple, amount int64) *UtxoSenderAmountTuple {
	return &UtxoSenderAmountTuple{
		sender: sender,
		amount: amount,
	}
}

type UtxoFromOutMap struct {
	mxp map[wire.OutPoint]*UtxoSenderAmountTuple
}

func NewUtxoFromOutMap(mxp map[wire.OutPoint]*UtxoSenderAmountTuple) *UtxoFromOutMap {
	return &UtxoFromOutMap{mxp: mxp}
}

func (uc UtxoFromOutMap) GetUtxoFrom(utxo wire.OutPoint) (*UtxoSenderAmountTuple, error) {
	utxoFrom, ok := uc.mxp[utxo]
	if !ok {
		return nil, errors.Errorf("not-exist-utxo[%s:%d]", utxo.Hash.String(), utxo.Index)
	}
	return utxoFrom, nil
}
