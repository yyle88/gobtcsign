package main

import (
	"fmt"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/yyle88/gobtcsign"
	"github.com/yyle88/gobtcsign/internal/utils"
)

func main() {
	const senderAddress = "tb1qvg2jksxckt96cdv9g8v9psreaggdzsrlm6arap"
	const privateKeyHex = "54bb1426611226077889d63c65f4f1fa212bcb42c2141c81e0c5409324711092" //注意不要暴露私钥，除非准备放弃这个钱包

	netParams := chaincfg.TestNet3Params

	param := gobtcsign.CustomParam{
		VinList: []gobtcsign.VinType{
			{
				OutPoint: *gobtcsign.MustNewOutPoint("e1f05d4ef10d6d4245839364c637cc37f429784883761668978645c67e723919", 2),
				Sender:   *gobtcsign.NewAddressTuple(senderAddress),
				Amount:   13089,
				RBFInfo:  *gobtcsign.NewRBFNotUse(),
			},
		},
		OutList: []gobtcsign.OutType{
			{ //转给第一个地址
				Target: *gobtcsign.NewAddressTuple("tb1qk0z8zhsq5hlewplv0039smnz62r2ujscz6gqjx"),
				Amount: 1234,
			},
			{ //注意，需要把剩下的转给自己。这是这笔交易的找零的输出，就是把剩下的钱再转给自己，但是也不要全部转给自己，还要留些就是矿工费用
				Target: *gobtcsign.NewAddressTuple(senderAddress),
				Amount: 11855 - 11111, //这里不要把剩下的钱都转给自己，要留一些矿工费 tx-fee，具体费用需要根据交易大小和费率计算的，这里随便给个假设的数量
			},
		},
		RBFInfo: *gobtcsign.NewRBFActive(),
	}

	//具体费用跟实时费率以及交易体大小有关，因此不同的交易有不同的预估值，这里省去预估的过程
	utils.MustEquals(int64(11111), int64(param.GetFee()))

	//得到待签名的交易
	signParam, err := param.GetSignParam(&netParams)
	utils.MustDone(err)

	fmt.Println(len(signParam.InputOuts))

	//签名
	utils.MustDone(gobtcsign.Sign(senderAddress, privateKeyHex, signParam))

	//这是签名后的交易
	msgTx := signParam.MsgTx

	//验证签名
	utils.MustDone(gobtcsign.VerifyP2PKHSignV2(msgTx, param.GetInputList(), &netParams))
	//比较信息
	utils.MustDone(gobtcsign.CheckMsgTxSameWithParam(msgTx, param, &netParams))

	//获得交易哈希
	txHash := gobtcsign.GetTxHash(msgTx)
	fmt.Println("msg-tx-hash:->", txHash, "<-")
	utils.MustEquals("e587e4f65a7fa5dbba6bede6b000e8ece097671bb348db3de0e507c8b36469ad", txHash)

	//把交易序列化得到hex字符串
	signedHex, err := gobtcsign.CvtMsgTxToHex(msgTx)
	utils.MustDone(err)
	fmt.Println("raw-tx-data:->", signedHex, "<-")
	utils.MustEquals("010000000001011939727ec645869768167683487829f437cc37c664938345426d0df14e5df0e10200000000fdffffff02d204000000000000160014b3c4715e00a5ff9707ec7be2586e62d286ae4a18e80200000000000016001462152b40d8b2cbac358541d850c079ea10d1407f02483045022100e8269080acc14fd24ee13cbbdaa5ea34192f090c917b4ca3da44eda25badd58e02206813da9023bebd556a95e04e6a55c9a5fdf5dfb19746c896d7fd7f26aaa58878012102407ea64d7a9e992028a94481af95ea7d8f54870bd73e5878a014da594335ba3200000000", signedHex)

	//SendRawHexTx(txHex) //通过这个tx-hex就可以发交易，我已经发完交易，你可以在链上看到它
	//假如报错 "-3: Amount is not a number or string" 说明发交易指令用的是 btcjson.NewSendRawTransactionCmd(txHex, &allowHighFees) 传的布尔值，而不是 btcjson.NewBitcoindSendRawTransactionCmd(txHex, maxFeeRate) 传的数字值。具体使用哪个主要是看节点的版本号，这里逻辑应当匹配
	//假如报错 "-26: mempool min fee not met, 2000 < 19940" 说明节点的 minrelaytxfee 设置的比较大，通常而言自建的测试节点的费用门槛要设置小些
	fmt.Println("success")
}

//当你发完交易以后查发送者的账户信息
//CONFIRMED UNSPENT	1 OUTPUTS (0.00013089 tBTC)
//UNCONFIRMED TX COUNT	1
//UNCONFIRMED RECEIVED	1 OUTPUTS (0.00000744 tBTC)
//UNCONFIRMED SPENT	1 OUTPUTS (0.00013089 tBTC)

//当你发完交易以后查接收者的账户信息
//CONFIRMED UNSPENT	1 OUTPUTS (0.00003000 tBTC)
//UNCONFIRMED TX COUNT	1
//UNCONFIRMED RECEIVED	1 OUTPUTS (0.00001234 tBTC)

//接下来等待链的确认即可，给的手续费越高确认越快，否则就需要耐心等待，或者提高手续费重新构造和发送交易
