package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestCustomParam_GetSignParam(t *testing.T) {
	//see dogecoin testnet hash: b80ced1c69ccbf16073e6ca48bc0a82c7c9bd8e08df21374dc771b9443ef6ac4

	//use testnet in this case
	netParams := dogecoin.TestNetParams

	//which tx the utxo from.
	utxoHash, err := chainhash.NewHashFromStr("8a0fb49fe4c407e24c8dd13e74a3398059ea3183082c0ea621a43d3500ee5918")
	require.NoError(t, err)

	//which address own the utxo. convert to pk-script bytes
	pkScript := caseGetAddressPkScript(t, "nXMSrjEQXUJ77TQSeErpJMySy3kfSfwSCP", netParams)

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
				Amount:  49921868563, //在已经确定utxo的来源hash和位置序号以后这里的数量其实是非必要的，但在某些格式的签名中是需要的
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

	res, err := customParam.GetSignParam(&netParams)
	require.NoError(t, err)
	require.Equal(t, res.NetParams.Net, netParams.Net)

	t.Log(len(res.PkScripts)) //只包含输入的数据
	require.Len(t, res.PkScripts, 1)
	require.Equal(t, pkScript, res.PkScripts[0])

	t.Log(res.Amounts) //只包含输入的数量
	require.Len(t, res.Amounts, 1)
	require.Equal(t, customParam.VinList[0].Amount, res.Amounts[0])

	require.Len(t, res.MsgTx.TxIn, 1)
	require.Equal(t, "8a0fb49fe4c407e24c8dd13e74a3398059ea3183082c0ea621a43d3500ee5918", res.MsgTx.TxIn[0].PreviousOutPoint.Hash.String())
	require.Equal(t, uint32(1), res.MsgTx.TxIn[0].PreviousOutPoint.Index)
	require.Nil(t, res.MsgTx.TxIn[0].SignatureScript) //由于还没有签名因此这里还是空的，将来签名是存在这里的，主要就是对输入的信息签名
	require.Equal(t, wire.MaxTxInSequenceNum-2, res.MsgTx.TxIn[0].Sequence)

	require.Len(t, res.MsgTx.TxOut, 2)
	require.Equal(t, int64(6547487), res.MsgTx.TxOut[0].Value)
	require.Equal(t, caseGetAddressPkScript(t, "nqNjvWut21qMKyZb4EPWBEUuVDSHuypVUa", netParams), res.MsgTx.TxOut[0].PkScript)
	require.Equal(t, int64(49914980576), res.MsgTx.TxOut[1].Value)
	require.Equal(t, caseGetAddressPkScript(t, "nXMSrjEQXUJ77TQSeErpJMySy3kfSfwSCP", netParams), res.MsgTx.TxOut[1].PkScript)
}
