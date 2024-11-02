package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestCustomParam_GetSignParam(t *testing.T) {
	//see dogecoin testnet hash: b80ced1c69ccbf16073e6ca48bc0a82c7c9bd8e08df21374dc771b9443ef6ac4

	//use testnet in this case
	netParams := dogecoin.TestNetParams

	const senderAddress = "nXMSrjEQXUJ77TQSeErpJMySy3kfSfwSCP"

	//which address own the utxo. convert to pk-script bytes
	pkScript := caseGetAddressPkScript(t, senderAddress, &netParams)

	customParam := CustomParam{
		VinList: []VinType{
			{
				OutPoint: *MustNewOutPoint(
					"8a0fb49fe4c407e24c8dd13e74a3398059ea3183082c0ea621a43d3500ee5918",
					1,
				),
				Sender: AddressTuple{ //这里是 address 或 pkScript 任填1个都行，这里示范填公钥脚本的情况
					PkScript: pkScript,
				},
				Amount:  49921868563,     //在已经确定utxo的来源hash和位置序号以后这里的数量其实是非必要的，但在某些格式的签名中是需要的
				RBFInfo: *NewRBFNotUse(), //这里不使用 RBF 机制，这个是控制单个utxo的
			},
		},
		OutList: []OutType{
			{
				//这里是 address 或 pkScript 任填1个都行，这里示范填地址的情况
				Target: *NewAddressTuple("nqNjvWut21qMKyZb4EPWBEUuVDSHuypVUa"),
				//发送数量
				Amount: 6547487,
			},
			{
				//这里是 address 或 pkScript 任填1个都行，这里示范填地址的情况
				Target: *NewAddressTuple(senderAddress),
				//找零数量
				Amount: 49914980576,
			},
		},
		RBFInfo: *NewRBFActive(), //这里使用 RBF 机制，这个是控制全部 utxo 的，优先级在单个utxo的后面
	}

	require.Equal(t, int64(340500), int64(customParam.GetFee()))

	res, err := customParam.GetSignParam(&netParams)
	require.NoError(t, err)
	require.Equal(t, res.NetParams.Net, netParams.Net)

	t.Log(len(res.InputOuts)) //只包含输入的数据
	require.Len(t, res.InputOuts, 1)
	require.Equal(t, pkScript, res.InputOuts[0].PkScript)

	t.Log(len(res.InputOuts)) //只包含输入的数量
	require.Len(t, res.InputOuts, 1)
	require.Equal(t, customParam.VinList[0].Amount, res.InputOuts[0].Value)

	require.Len(t, res.MsgTx.TxIn, 1)
	require.Equal(t, "8a0fb49fe4c407e24c8dd13e74a3398059ea3183082c0ea621a43d3500ee5918", res.MsgTx.TxIn[0].PreviousOutPoint.Hash.String())
	require.Equal(t, uint32(1), res.MsgTx.TxIn[0].PreviousOutPoint.Index)
	require.Nil(t, res.MsgTx.TxIn[0].SignatureScript) //由于还没有签名因此这里还是空的，将来签名是存在这里的，主要就是对输入的信息签名
	require.Equal(t, wire.MaxTxInSequenceNum-2, res.MsgTx.TxIn[0].Sequence)

	require.Len(t, res.MsgTx.TxOut, 2)
	require.Equal(t, int64(6547487), res.MsgTx.TxOut[0].Value)
	require.Equal(t, caseGetAddressPkScript(t, "nqNjvWut21qMKyZb4EPWBEUuVDSHuypVUa", &netParams), res.MsgTx.TxOut[0].PkScript)
	require.Equal(t, int64(49914980576), res.MsgTx.TxOut[1].Value)
	require.Equal(t, caseGetAddressPkScript(t, senderAddress, &netParams), res.MsgTx.TxOut[1].PkScript)
}
