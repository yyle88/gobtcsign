package dogecoin

import "github.com/yyle88/gobtcsign/internal/dusts"

const (
	// MinDustOutput 硬性灰尘数量，详见 https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md
	MinDustOutput = 100000 // The hard dust limit is set at 0.001 DOGE - outputs under this value are invalid and rejected.

	// SoftDustLimit 弹性灰尘限制，详见 https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md
	SoftDustLimit = 1000000 // The soft dust limit is set at 0.01 DOGE - sending a transaction with outputs under this value, are required to add 0.01 DOGE for each such output

	// ExtraDustsFee 额外的灰尘费，这会让手续费变得不稳定，让代码中所有 txrules.FeeForSerializeSize 的地方都附带额外的灰尘费
	ExtraDustsFee = 1000000 // add 0.01 DOGE for each such output
)

type DustFee = dusts.DustFee

func NewDogeDustFee() DustFee {
	res := dusts.NewDustFee()
	res.SoftDustSize = SoftDustLimit
	res.ExtraDustFee = ExtraDustsFee
	return res
}
