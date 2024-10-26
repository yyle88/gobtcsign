package gobtcsign

import (
	"encoding/base64"
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestCustomParam_GetSignParam(t *testing.T) {
	//see dogecoin testnet hash: b80ced1c69ccbf16073e6ca48bc0a82c7c9bd8e08df21374dc771b9443ef6ac4

	netParams := dogecoin.TestNetParams

	utxoHash, err := chainhash.NewHashFromStr("8a0fb49fe4c407e24c8dd13e74a3398059ea3183082c0ea621a43d3500ee5918")
	require.NoError(t, err)

	address, err := btcutil.DecodeAddress("nXMSrjEQXUJ77TQSeErpJMySy3kfSfwSCP", &netParams)
	require.NoError(t, err)

	pkScript, err := txscript.PayToAddrScript(address)
	require.NoError(t, err)

	{ // 这里写个简单的比较逻辑
		expectedBytes, err := base64.StdEncoding.DecodeString("dqkUIqn5GrQ9r5dmyzOQ/RTb+qkqVQqIrA==")
		require.NoError(t, err)
		require.Equal(t, expectedBytes, pkScript)
	}

	customParam := CustomParam{
		VinList: []VinType{
			{
				OutPoint: *wire.NewOutPoint(
					utxoHash, //这个是收到utxo的交易哈希，即utxo是从哪里来的，配合位置索引序号构成唯一索引，就能确定是花的哪个utxo
					1,        //这个是收到utxo的位置，比如一个交易中有多个输出，这里要选择输出的位置
				),
				Sender: AddressTuple{ //这里是 address 或 pkScript 任填1个都行
					PkScript: pkScript,
				},
				Amount:  49921868563,
				RBFInfo: RBFConfig{}, //这里不使用 RBF 机制，这个是控制单个utxo的
			},
		},
		OutList: []OutType{
			{
				Target: AddressTuple{ //这里是 address 或 pkScript 任填1个都行
					Address: "nqNjvWut21qMKyZb4EPWBEUuVDSHuypVUa",
				},
				Amount: 6547487, //发送数量
			},
			{
				Target: AddressTuple{ //这里是 address 或 pkScript 任填1个都行
					Address: "nXMSrjEQXUJ77TQSeErpJMySy3kfSfwSCP",
				},
				Amount: 49914980576, //找零数量
			},
		},
		RBFInfo: RBFConfig{ //这里使用 RBF 机制，这个是控制全部 utxo 的，优先级在单个utxo的后面
			AllowRBF: true,
			Sequence: wire.MaxTxInSequenceNum - 2, //默认的RBF机制就是这样的
		},
	}

	signParam, err := customParam.GetSignParam(&netParams)
	require.NoError(t, err)

	t.Log(len(signParam.PkScripts)) //只包含输入的数据
	require.Len(t, signParam.PkScripts, 1)
	require.Equal(t, pkScript, signParam.PkScripts[0])

	t.Log(signParam.Amounts) //只包含输入的数量
	require.Len(t, signParam.Amounts, 1)
	require.Equal(t, customParam.VinList[0].Amount, signParam.Amounts[0])

	require.Equal(t, signParam.NetParams.Net, netParams.Net)
}
