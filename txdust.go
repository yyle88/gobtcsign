package gobtcsign

import (
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/yyle88/gobtcsign/internal/dusts"
)

type DustFee = dusts.DustFee

func NewDustFee() DustFee {
	return dusts.NewDustFee() //比特币没有软灰尘收费，这里配置个空的（因为doge里有，这里为了逻辑相通，而给个空的）
}

type DustLimit = dusts.DustLimit

func NewDustLimit() *DustLimit {
	return dusts.NewDustLimit(txrules.IsDustOutput)
}
