package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestIsDustOutput_BTC(t *testing.T) {
	netParams := chaincfg.MainNetParams

	dustLimit := NewDustLimit()

	const address = "1MAgFFbMpgx6hTPvK3HY348Bbnwk6RFHm5"
	const feeRate = 1234000 //假设手续费是这个数

	amount := int64(100000)
	for {
		output := wire.NewTxOut(amount, MustGetPkScript(MustNewAddress(address, &netParams)))
		if !dustLimit.IsDustOutput(output, feeRate) {
			break
		}
		amount++
	}
	t.Log("amount:", amount, "IS NOT DUST IN BTC")
	t.Log("amount:", amount-1, "IS DUST IN BTC")
}

func TestIsDustOutput_BTC_2(t *testing.T) {
	netParams := chaincfg.MainNetParams

	dustLimit := NewDustLimit()

	const address = "1MAgFFbMpgx6hTPvK3HY348Bbnwk6RFHm5"
	const feeRate = 1000 //假设手续费是这个数

	amount := int64(1)
	for {
		output := wire.NewTxOut(amount, MustGetPkScript(MustNewAddress(address, &netParams)))
		if !dustLimit.IsDustOutput(output, feeRate) {
			break
		}
		amount++
	}
	require.Equal(t, int64(546), amount) //这就是有的教程和代码里写 546 的依据
	// 546 是比特币网络中一个常见的 Dust Threshold（灰尘阈值），其来源与比特币的交易手续费和 UTXO（未花费交易输出，Unspent Transaction Output）管理策略相关。
	// 具体来说，546 是在默认设置下，比特币 Core 客户端用于计算 Dust Output（灰尘输出）的默认值。
	// 比特币网络中最早的版本采用的是类似值，后来进行了微调，但很多第三方实现（例如钱包）仍然保留 546 的惯例。
	t.Log("amount:", amount, "IS NOT DUST IN BTC")
	t.Log("amount:", amount-1, "IS DUST IN BTC")
}

func TestIsDustOutput_BTC_3(t *testing.T) {
	netParams := chaincfg.TestNet3Params

	dustLimit := NewDustLimit()

	const address = "tb1qy2f7svy0hp57wz3p6hvu0vf5fys750932ct3q5"
	const feeRate = 1000 //假设手续费是这个数

	amount := int64(1)
	for {
		output := wire.NewTxOut(amount, MustGetPkScript(MustNewAddress(address, &netParams)))
		if !dustLimit.IsDustOutput(output, feeRate) {
			break
		}
		amount++
	}
	require.Equal(t, int64(294), amount) //当地址类型为 P2WPKH 时由于其 输入大小比 P2PKH 小得多，因此灰尘阈值也更小些
	// 因此有的代码里会这样写
	// DustLimit returns the output dust limit (lowest possible satoshis in a UTXO) for the address type.
	// func (a Address) DustLimit() int64 {
	//		switch a.encodedType {
	//		case AddressP2TR:
	//			return 330
	//		case AddressP2WPKH:
	//			return 294
	//		default:
	//			return 546
	//		}
	// }
	// 但实际上大家还都是使用 out >= 546 作为限制
	t.Log("amount:", amount, "IS NOT DUST IN BTC")
	t.Log("amount:", amount-1, "IS DUST IN BTC")
}

func TestIsDustOutput_DOGE(t *testing.T) {
	netParams := dogecoin.TestNetParams

	dustLimit := dogecoin.NewDogeDustLimit()
	{
		const amount = 500
		const address = "nr2XmwqixAdXwkgVyshx3HPFRMfXugM8Zi"
		const feeRate = 0 //这个不重要

		output := wire.NewTxOut(amount, MustGetPkScript(MustNewAddress(address, &netParams)))
		require.True(t, dustLimit.IsDustOutput(output, feeRate))
		t.Log("amount:", amount, "IS DUST IN DOGE")
	}
	{
		const amount = 100000
		const address = "nqedQEDCgwrXqLd2JrrpCfD9Tcz384rdHA"
		const feeRate = 0 //这个不重要

		output := wire.NewTxOut(amount, MustGetPkScript(MustNewAddress(address, &netParams)))
		require.False(t, dustLimit.IsDustOutput(output, feeRate))
		t.Log("amount:", amount, "IS NOT DUST IN DOGE")
	}
}
