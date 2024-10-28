package gobtcsign

import (
	"encoding/base64"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
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

func TestGetAddressPkScript(t *testing.T) {
	netParams := dogecoin.TestNetParams
	pkScript := caseGetAddressPkScript(t, "nXMSrjEQXUJ77TQSeErpJMySy3kfSfwSCP", netParams)

	// 这里写个简单的比较逻辑
	pkTarget, err := base64.StdEncoding.DecodeString("dqkUIqn5GrQ9r5dmyzOQ/RTb+qkqVQqIrA==")
	require.NoError(t, err)
	require.Equal(t, pkTarget, pkScript)
}

func caseGetAddressPkScript(t *testing.T, rawAddress string, netParams chaincfg.Params) []byte {
	pkScript, err := GetAddressPkScript(rawAddress, &netParams)
	require.NoError(t, err)
	return pkScript
}
