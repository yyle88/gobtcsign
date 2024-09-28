package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/yyle88/gobtcsign/dogecoin"
)

type DustFeeConfig struct {
	SoftDustLimit btcutil.Amount
	ExtraDustsFee btcutil.Amount
}

func NewDustFeeConfig() DustFeeConfig {
	return DustFeeConfig{
		SoftDustLimit: 0,
		ExtraDustsFee: 0,
	}
}

func NewDogeDustFeeConfig() DustFeeConfig {
	return DustFeeConfig{
		SoftDustLimit: dogecoin.ChainSoftDustLimit,
		ExtraDustsFee: dogecoin.ChainExtraDustsFee,
	}
}

func (dfc *DustFeeConfig) GetSoftDustSize(txOuts []*wire.TxOut) int64 {
	var amount = int64(dfc.SoftDustLimit)
	if amount == 0 {
		return 0
	}

	var n int64
	for _, x := range txOuts {
		if x.Value < amount {
			n++
		}
	}
	return n
}

func (dfc *DustFeeConfig) GetSoftDustsFee(txOuts []*wire.TxOut) btcutil.Amount {
	return btcutil.Amount(dfc.GetSoftDustSize(txOuts) * int64(dfc.ExtraDustsFee))
}
