package dusts

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

type DustFee struct {
	SoftDustSize btcutil.Amount //这是软灰尘限制，即介于软灰尘和硬灰尘之间的数量
	ExtraDustFee btcutil.Amount //单个软灰尘额外的收费
}

func NewDustFee() DustFee {
	return DustFee{
		SoftDustSize: 0,
		ExtraDustFee: 0,
	}
}

func (D *DustFee) CountDustOutput(outputs []*wire.TxOut) int64 {
	var minLimit = int64(D.SoftDustSize)
	if minLimit == 0 {
		return 0
	}

	var count int64
	for _, out := range outputs {
		if out.Value < minLimit {
			count++
		}
	}
	return count
}

func (D *DustFee) SumExtraDustFee(outputs []*wire.TxOut) btcutil.Amount {
	return btcutil.Amount(D.CountDustOutput(outputs) * int64(D.ExtraDustFee))
}
