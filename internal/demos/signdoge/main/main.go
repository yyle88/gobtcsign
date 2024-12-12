package main

import (
	"fmt"

	"github.com/yyle88/gobtcsign"
	"github.com/yyle88/gobtcsign/dogecoin"
	"github.com/yyle88/gobtcsign/internal/utils"
)

func main() {
	const senderAddress = "nkgVWbNrUowCG4mkWSzA7HHUDe3XyL2NaC"
	const privateKeyHex = "5f397bc72377b75db7b008a9c3fcd71651bfb138d6fc2458bb0279b9cfc8442a" //注意不要暴露私钥，除非准备放弃这个钱包

	netParams := dogecoin.TestNetParams

	param := gobtcsign.BitcoinTxParams{
		VinList: []gobtcsign.VinType{
			{
				OutPoint: *gobtcsign.MustNewOutPoint(
					"173d5e1b33fc9adf64cd4b1f3b2ac73acaf0e10c967cd6fa1aa191d817d7ff77",
					3, //这里的位置是3，哈希和位置构成UTXO的主键，这里能从链上查到位置
				),
				Sender:  *gobtcsign.NewAddressTuple(senderAddress),
				Amount:  14049272,
				RBFInfo: *gobtcsign.NewRBFNotUse(),
			},
		},
		OutList: []gobtcsign.OutType{
			{
				Target: *gobtcsign.NewAddressTuple("ng4P16anXNUrQw6VKHmoMW8NHsTkFBdNrn"),
				Amount: 1234567,
			},
			{ //注意，需要把剩下的转给自己。这是这笔交易的找零的输出，就是把剩下的钱再转给自己，但是也不要全部转给自己，还要留些就是矿工费用
				Target: *gobtcsign.NewAddressTuple(senderAddress),
				Amount: 12814705 - 222222, //这里不要把剩下的钱都转给自己，要留一些矿工费 tx-fee，具体费用需要根据交易大小和费率计算的，这里随便给个假设的数量
			},
		},
		RBFInfo: *gobtcsign.NewRBFActive(),
	}

	//具体费用跟实时费率以及交易体大小有关，因此不同的交易有不同的预估值，这里省去预估的过程
	utils.MustEquals(int64(222222), int64(param.GetFee()))

	size, err := param.EstimateTxSize(&netParams, gobtcsign.NewNoChange())
	utils.MustDone(err)
	fmt.Println("estimate-tx-size:", size) //这是预估值 略微 >= 实际值

	//得到待签名的交易
	signParam, err := param.CreateTxSignParams(&netParams)
	utils.MustDone(err)

	//签名
	utils.MustDone(gobtcsign.Sign(senderAddress, privateKeyHex, signParam))

	//这是签名后的交易
	msgTx := signParam.MsgTx

	//验证签名
	utils.MustDone(param.VerifyMsgTxSign(msgTx, &netParams))
	//比较信息
	utils.MustDone(param.CheckMsgTxParam(msgTx, &netParams))

	//获得交易哈希
	txHash := gobtcsign.GetTxHash(msgTx)
	fmt.Println("msg-tx-hash:->", txHash, "<-")
	utils.MustEquals("d06f0a49c4f18e2aa520eb3bfc961602aa18c811380cb38cae3638c13883f5ed", txHash)

	//把交易序列化得到hex字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	utils.MustDone(err)
	fmt.Println("raw-tx-data:->", signedHex, "<-")
	utils.MustEquals("010000000177ffd717d891a11afad67c960ce1f0ca3ac72a3b1f4bcd64df9afc331b5e3d17030000006a473044022025a41ebdb7d1a5edc5bcdb120ac339591fd95a9a084c8250a362073ffb27575202204579fa82476a52f5a28f605a827ef4866d4ba671c60363f22b523f5c27bf090a012102dfef3896f159dde1c2a972038e06ebc39c551f5f3d45e2fc9544f951fe4282f4fdffffff0287d61200000000001976a9148228d0af289894d419ddcaf6da679d8e9f0f160188ac6325c000000000001976a914b4ddb9db68061a0fec90a4bcaef21f82c8cfa1eb88ac00000000", signedHex)

	//SendRawHexTx(txHex) //通过这个tx-hex就可以发交易，我已经发完交易，你可以在链上看到它
	fmt.Println("success")
}
