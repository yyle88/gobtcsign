package dusts

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
)

type DustLimit struct {
	check func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool
}

func NewDustLimit(check func(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool) *DustLimit {
	return &DustLimit{check: check}
}

func (D *DustLimit) IsDustOutput(output *wire.TxOut, relayFeePerKb btcutil.Amount) bool {
	return D.check(output, relayFeePerKb)
}
