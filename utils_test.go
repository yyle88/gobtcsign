package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

func TestIsDustOutputDoge(t *testing.T) {
	const dustLimitCoin = 0.001 //THE HARD DUST LIMIT IS SET AT 0.001 DOGE - OUTPUTS UNDER THIS VALUE ARE INVALID AND REJECTED. (https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md)
	{
		const amount = 500
		isDust, err := IsDustOutputDoge(amount, dustLimitCoin)
		require.NoError(t, err)
		require.True(t, isDust)
		t.Log("amount:", amount, "IS DUST IN DOGE")
	}
	{
		const amount = 100000
		isDust, err := IsDustOutputDoge(amount, dustLimitCoin)
		require.NoError(t, err)
		require.False(t, isDust)
		t.Log("amount:", amount, "IS NOT DUST IN DOGE")
	}
}

func TestIsDustOutputBtc(t *testing.T) {
	amount := int64(100000)
	for {
		const address = "1MAgFFbMpgx6hTPvK3HY348Bbnwk6RFHm5"
		const feePerKb = 1234 * 1000 //假设手续费是这个数
		const dustLimitCoin = 546e-8 //BTC CHAIN: CONSIDERS ANYTHING BELOW 546 SATOSHIS (PARTS OF A BITCOIN) TO BE DUST.
		isDust, err := IsDustOutputBtc(&chaincfg.MainNetParams, address, amount, dustLimitCoin, feePerKb)
		require.NoError(t, err)
		if !isDust {
			break
		}
		amount++
	}
	t.Log("amount:", amount, "IS NOT DUST IN BTC")
	t.Log("amount:", amount-1, "IS DUST IN BTC")
}
