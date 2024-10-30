package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

type DustFee struct {
	SoftDustLimit btcutil.Amount
	ExtraDustsFee btcutil.Amount
}

func NewDustFee() DustFee {
	return DustFee{
		SoftDustLimit: 0,
		ExtraDustsFee: 0,
	}
}

func (D *DustFee) CountDustsOutNum(txOuts []*wire.TxOut) int64 {
	var dustLimit = int64(D.SoftDustLimit)
	if dustLimit == 0 {
		return 0
	}

	var count int64
	for _, out := range txOuts {
		if out.Value < dustLimit {
			count++
		}
	}
	return count
}

func (D *DustFee) SumExtraDustsFee(txOuts []*wire.TxOut) btcutil.Amount {
	return btcutil.Amount(D.CountDustsOutNum(txOuts) * int64(D.ExtraDustsFee))
}
