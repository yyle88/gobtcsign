package dogecoin

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/yyle88/gobtcsign/internal/dusts"
)

const (
	// MinDustOutput 硬性灰尘数量，详见 https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md
	MinDustOutput = 100000 // The hard dust limit is set at 0.001 DOGE - outputs under this value are invalid and rejected.

	// SoftDustLimit 弹性灰尘限制，详见 https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md
	SoftDustLimit = 1000000 // The soft dust limit is set at 0.01 DOGE - sending a transaction with outputs under this value, are required to add 0.01 DOGE for each such output

	// ExtraDustsFee 额外的灰尘费，这会让手续费变得不稳定，让代码中所有 txrules.FeeForSerializeSize 的地方都附带额外的灰尘费
	ExtraDustsFee = 1000000 // add 0.01 DOGE for each such output
)

type DustFee = dusts.DustFee

// NewDogeDustFee 配置狗狗币的dust费用
// 具体参考链接在
// https://github.com/dogecoin/dogecoin/blob/b4a5d2bef20f5cca54d9c14ca118dec259e47bb4/doc/fee-recommendation.md
// DOGECOIN 简单规定了软灰尘和硬灰尘，假如是硬灰尘会被拒绝，假如是软灰尘会收取额外的费用
func NewDogeDustFee() DustFee {
	res := dusts.NewDustFee()
	res.SoftDustSize = SoftDustLimit
	res.ExtraDustFee = ExtraDustsFee
	return res
}

type DustLimit = dusts.DustLimit

func NewDogeDustLimit() *DustLimit {
	return dusts.NewDustLimit(func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool {
		return output.Value < MinDustOutput //在dogecoin中的灰尘规定比较简单，它不依赖于费率，而是直接和常量比较，逻辑简单
	})
}
