package gobtcsign_test

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestSignBTC(t *testing.T) {
	const senderAddress = "tb1qvg2jksxckt96cdv9g8v9psreaggdzsrlm6arap"
	const privateKeyHex = "54bb1426611226077889d63c65f4f1fa212bcb42c2141c81e0c5409324711092" //注意不要暴露私钥，除非准备放弃这个钱包

	netParams := chaincfg.TestNet3Params

	param := &gobtcsign.CustomParam{
		VinList: []gobtcsign.VinType{
			{
				OutPoint: *gobtcsign.MustNewOutPoint("fb87cc4010bd4a34cb4be86f37182fada63c9923ae8eae5d2f793cb5f50c6328", 0),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   4900,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
			{
				OutPoint: *gobtcsign.MustNewOutPoint("fcc889d7f0217694ab46d93f03a200d326c34e317552a6a33cb3fab03aa0b439", 1),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   4320,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
			{
				OutPoint: *gobtcsign.MustNewOutPoint("5c98431bbb271ea3652168d2b4da8a76573fd8fec104e73f6f6f3a7c6fe6b97d", 0),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   4560,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
			{
				OutPoint: *gobtcsign.MustNewOutPoint("5fe7486105cb41cc1496fed89296140e00fee5fdc880ac335ea1df9b374f9348", 3),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   4900,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
			{
				OutPoint: *gobtcsign.MustNewOutPoint("de8ac7275793df0218d7151e420393fa3cf39159147fa6453c5f279f249d6a52", 1),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   22865,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
		},
		OutList: []gobtcsign.OutType{
			{ //转给第一个地址
				Target: *gobtcsign.NewAddressTuple("tb1qk0z8zhsq5hlewplv0039smnz62r2ujscz6gqjx"),
				Amount: 3000,
			},
			{ //转给第二个地址
				Target: *gobtcsign.NewAddressTuple("tb1qlj64u6fqutr0xue85kl55fx0gt4m4urun25p7q"),
				Amount: 2000,
			},
			{ //注意，需要把剩下的转给自己。这是这笔交易的找零的输出，就是把剩下的钱再转给自己，但是也不要全部转给自己，还要留些就是矿工费用
				Target: *gobtcsign.NewAddressTuple(senderAddress),
				Amount: 36545 - 23456, //这里不要把剩下的钱都转给自己，要留一些矿工费 tx-fee，具体费用需要根据交易大小和费率计算的，这里随便给个假设的数量
			},
		},
		RBFInfo: *gobtcsign.NewRBFActive(),
	}

	//具体费用跟实时费率以及交易体大小有关，因此不同的交易有不同的预估值，这里省去预估的过程
	require.Equal(t, int64(23456), int64(param.GetFee()))

	//得到待签名的交易
	signParam, err := param.GetSignParam(&netParams)
	require.NoError(t, err)

	t.Log(len(signParam.InputOuts))

	//签名
	require.NoError(t, gobtcsign.Sign(senderAddress, privateKeyHex, signParam))

	//这是签名后的交易
	msgTx := signParam.MsgTx

	//验证签名
	require.NoError(t, gobtcsign.VerifySignV2(msgTx, param.GetInputList(), &netParams))
	//比较信息
	require.NoError(t, gobtcsign.CheckMsgTxSameWithParam(msgTx, *param, &netParams))

	//获得交易哈希
	txHash := gobtcsign.GetTxHash(msgTx)
	t.Log("msg-tx-hash:->", txHash, "<-")
	require.Equal(t, "e1f05d4ef10d6d4245839364c637cc37f429784883761668978645c67e723919", txHash)

	//把交易序列化得到hex字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	require.NoError(t, err)
	t.Log("raw-tx-data:->", signedHex, "<-")
	require.Equal(t, "0100000000010528630cf5b53c792f5dae8eae23993ca6ad2f18376fe84bcb344abd1040cc87fb0000000000fdffffff39b4a03ab0fab33ca3a65275314ec326d300a2033fd946ab947621f0d789c8fc0100000000fdffffff7db9e66f7c3a6f6f3fe704c1fed83f57768adab4d2682165a31e27bb1b43985c0000000000fdffffff48934f379bdfa15e33ac80c8fde5fe000e149692d8fe9614cc41cb056148e75f0300000000fdffffff526a9d249f275f3c45a67f145991f33cfa9303421e15d71802df935727c78ade0100000000fdffffff03b80b000000000000160014b3c4715e00a5ff9707ec7be2586e62d286ae4a18d007000000000000160014fcb55e6920e2c6f37327a5bf4a24cf42ebbaf07c213300000000000016001462152b40d8b2cbac358541d850c079ea10d1407f0247304402201977b7da04ca36eb2b7cb3334b7f0f6551733ee02097c672f24938766b47778e02202bc7b30a01b55748f408f35a666376ee23695490d349cde54fee7a00bf5e6926012102407ea64d7a9e992028a94481af95ea7d8f54870bd73e5878a014da594335ba3202483045022100d06785151683bf48e2327194ba0a09090e76d9e7f1d1d2f67c69b693c3663127022001f553f91afc9820507c32f1a0197719df2312e51a80fc3d517efcc2bc14168c012102407ea64d7a9e992028a94481af95ea7d8f54870bd73e5878a014da594335ba3202483045022100f7aee523be19affa0548dcd3f74352088da60d2630a4cbac55f81757a5de112f02203e9a09a2b90df4ab61a440a3b72b4448a250289c0680db992ade9080a19cdc44012102407ea64d7a9e992028a94481af95ea7d8f54870bd73e5878a014da594335ba3202473044022079092a4916cac06227c914df0baeca915f2f5f7c62983287a8238f933fc47f70022008140a0e7f1353461ad471a17aa405f8ba18bc028772896d61914b6d2bec6444012102407ea64d7a9e992028a94481af95ea7d8f54870bd73e5878a014da594335ba320247304402205396b130a34f9bd361cc9b150b3053a31b8225234685f5585729e26f1ddd4c540220349d16280aec9e1ca8d18a556f5e9cb740f18362d36337bc81c3d23cf8ff3b1c012102407ea64d7a9e992028a94481af95ea7d8f54870bd73e5878a014da594335ba3200000000", signedHex)

	//SendRawHexTx(txHex) //通过这个tx-hex就可以发交易，我已经发完交易，你可以在链上看到它
	//假如报错 "-3: Amount is not a number or string" 说明发交易指令用的是 btcjson.NewSendRawTransactionCmd(txHex, &allowHighFees) 传的布尔值，而不是 btcjson.NewBitcoindSendRawTransactionCmd(txHex, maxFeeRate) 传的数字值。具体使用哪个主要是看节点的版本号，这里逻辑应当匹配
	//假如报错 "-26: mempool min fee not met, 2000 < 19940" 说明节点的 minrelaytxfee 设置的比较大，通常而言自建的测试节点的费用门槛要设置小些
	t.Log("success")
}

func TestSignDOGE(t *testing.T) {
	const senderAddress = "nkgVWbNrUowCG4mkWSzA7HHUDe3XyL2NaC"
	const privateKeyHex = "5f397bc72377b75db7b008a9c3fcd71651bfb138d6fc2458bb0279b9cfc8442a" //注意不要暴露私钥，除非准备放弃这个钱包

	netParams := dogecoin.TestNetParams

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
	signParam, err := param.GetSignParam(&netParams)
	require.NoError(t, err)

	//签名
	require.NoError(t, gobtcsign.Sign(senderAddress, privateKeyHex, signParam))

	//这是签名后的交易
	msgTx := signParam.MsgTx

	//验证签名
	require.NoError(t, gobtcsign.VerifySignV2(msgTx, param.GetInputList(), &netParams))
	//比较信息
	require.NoError(t, gobtcsign.CheckMsgTxSameWithParam(msgTx, param, &netParams))

	//获得交易哈希
	txHash := gobtcsign.GetTxHash(msgTx)
	t.Log("msg-tx-hash:->", txHash, "<-")
	require.Equal(t, "173d5e1b33fc9adf64cd4b1f3b2ac73acaf0e10c967cd6fa1aa191d817d7ff77", txHash)

	//把交易序列化得到hex字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	require.NoError(t, err)
	t.Log("raw-tx-data:->", signedHex, "<-")
	require.Equal(t, "01000000020befa5cd809bb70ccd0d435106c1c42849af4a207092d4cbc5f4d3654851a357000000006b4830450221009295278402e377ec62ad54c27ffb894960e4cf8935cddd968824f1aab62c1e2b02206d71cc98b48a0513ddf3f04271a989a27aabda333afe51323348cf9ec5cfe23f012102dfef3896f159dde1c2a972038e06ebc39c551f5f3d45e2fc9544f951fe4282f4fdffffffb48c8b92824e0d25d65dfd7ad7fff73d046f74b811e86fbc40591c2289c93eaf000000006a473044022066f8c2dc9387d7627bcf52feeb5b02e45b8bf0228c3bdc065844dadef5f2708e02206a6e2b809b90117e68fa0e4ca3381ccf0c9f0f770231dbc45f8dcb0f7e2cd58f012102dfef3896f159dde1c2a972038e06ebc39c551f5f3d45e2fc9544f951fe4282f4fdffffff0487d61200000000001976a9148228d0af289894d419ddcaf6da679d8e9f0f160188acceca2300000000001976a914633a7a97acb866a45cf4e77bb8527f3e8bc2bdd788ac15bf3400000000001976a914c58ad3d72c51f2fcf83fe10d70de8a91cf11dba988acf85fd600000000001976a914b4ddb9db68061a0fec90a4bcaef21f82c8cfa1eb88ac00000000", signedHex)

	//SendRawHexTx(txHex) //通过这个tx-hex就可以发交易，我已经发完交易，你可以在链上看到它
	t.Log("success")
}
