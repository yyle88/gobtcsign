package gobtcsign_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestSignBTC(t *testing.T) {
	// can sign Bitcoin transaction
	t.Log("请保护好自己的私钥")
}

func TestSignDOGE(t *testing.T) {
	const senderAddress = "nkgVWbNrUowCG4mkWSzA7HHUDe3XyL2NaC"
	const privateKeyHex = "5f397bc72377b75db7b008a9c3fcd71651bfb138d6fc2458bb0279b9cfc8442a" //注意不要暴露私钥，除非准备放弃这个钱包

	param := gobtcsign.CustomParam{
		VinList: []gobtcsign.VinType{
			{
				OutPoint: *gobtcsign.MustNewOutPoint("57a3514865d3f4c5cbd49270204aaf4928c4c10651430dcd0cb79b80cda5ef0b", 0),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   6799372,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
			{
				OutPoint: *gobtcsign.MustNewOutPoint("af3ec989221c5940bc6fe811b8746f043df7ffd77afd5dd6250d4e82928b8cb4", 0),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   14632612,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
		},
		OutList: []gobtcsign.OutType{
			{
				Target: *gobtcsign.NewAddressTuple("ng4P16anXNUrQw6VKHmoMW8NHsTkFBdNrn"),
				Amount: 1234567,
			},
			{
				Target: *gobtcsign.NewAddressTuple("ndEqDSpcZquZspz5uro1M21ENmrz9Gbp9K"),
				Amount: 2345678,
			},
			{
				Target: *gobtcsign.NewAddressTuple("nnCfgxxyuJvYDanY3Y2nQMxW9wWWfGjkvR"),
				Amount: 3456789,
			},
			{ //注意，需要把剩下的转给自己。这是这笔交易的找零的输出，就是把剩下的钱再转给自己，但是也不要全部转给自己，还要留些就是矿工费用
				Target: *gobtcsign.NewAddressTuple(senderAddress),
				Amount: 14394950 - 345678, //这里不要把剩下的钱都转给自己，要留一些矿工费 tx-fee，具体费用需要根据交易大小和费率计算的，这里随便给个假设的数量
			},
		},
		RBFInfo: *gobtcsign.NewRBFActive(),
	}

	//具体费用跟实时费率以及交易体大小有关，因此不同的交易有不同的预估值，这里省去预估的过程
	require.Equal(t, int64(345678), int64(param.GetFee()))

	//得到待签名的交易
	signParam, err := param.GetSignParam(&dogecoin.TestNetParams)
	require.NoError(t, err)

	//签名
	require.NoError(t, gobtcsign.Sign(senderAddress, privateKeyHex, signParam))

	//这是签名后的交易
	msgTx := signParam.MsgTx

	//验证签名
	require.NoError(t, gobtcsign.VerifyP2PKHSignV2(msgTx, param.GetInputList(), &dogecoin.TestNetParams))
	//比较信息
	require.NoError(t, gobtcsign.CheckMsgTxSameWithParam(msgTx, param, &dogecoin.TestNetParams))

	//获得交易哈希
	txHash := gobtcsign.GetTxHash(msgTx)
	t.Log("msg-tx-hash:->", txHash, "<-")
	require.Equal(t, "173d5e1b33fc9adf64cd4b1f3b2ac73acaf0e10c967cd6fa1aa191d817d7ff77", txHash)

	//把交易序列化得到hex字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	require.NoError(t, err)
	t.Log("raw-tx-data:->", signedHex, "<-")
	require.Equal(t, "01000000020befa5cd809bb70ccd0d435106c1c42849af4a207092d4cbc5f4d3654851a357000000006b4830450221009295278402e377ec62ad54c27ffb894960e4cf8935cddd968824f1aab62c1e2b02206d71cc98b48a0513ddf3f04271a989a27aabda333afe51323348cf9ec5cfe23f012102dfef3896f159dde1c2a972038e06ebc39c551f5f3d45e2fc9544f951fe4282f4fdffffffb48c8b92824e0d25d65dfd7ad7fff73d046f74b811e86fbc40591c2289c93eaf000000006a473044022066f8c2dc9387d7627bcf52feeb5b02e45b8bf0228c3bdc065844dadef5f2708e02206a6e2b809b90117e68fa0e4ca3381ccf0c9f0f770231dbc45f8dcb0f7e2cd58f012102dfef3896f159dde1c2a972038e06ebc39c551f5f3d45e2fc9544f951fe4282f4fdffffff0487d61200000000001976a9148228d0af289894d419ddcaf6da679d8e9f0f160188acceca2300000000001976a914633a7a97acb866a45cf4e77bb8527f3e8bc2bdd788ac15bf3400000000001976a914c58ad3d72c51f2fcf83fe10d70de8a91cf11dba988acf85fd600000000001976a914b4ddb9db68061a0fec90a4bcaef21f82c8cfa1eb88ac00000000", signedHex)

	//SendRawHexTx(txHex) //通过这个tx-hex就可以发交易啦，我已经发完交易，你可以在链上看到它
	t.Log("success")
}
