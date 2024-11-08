package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
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

	customParam := &CustomParam{
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

func TestCustomParam_VerifyMsgTxSign(t *testing.T) {
	const txHex = "02000000000101e215fce89be7be2077954a8715d4257aef54750e8099e745a2c863eb50446ba00100000000fdffffff02ee5deb0200000000160014aa9b39359ea7a4146d8805d7507815a725c7bb86e8030000000000001600144ec2a17c77a078dd8a5eaef4ff4f4c15d522b14702473044022017c75d43739320246cb1f82d3f7c501aa0591849364781ededc9dd757e84636c022064e15135fab0d6bae0fd728cad6e1a6fd58a5da401d358cbd71b6798d30cc511012103840eed538d8060fdf812486faddb29ab18401a0f099ae19bbc9f23b80ca465f740c53000"

	msgTx, err := NewMsgTxFromHex(txHex)
	require.NoError(t, err)
	t.Log(msgTx.TxHash().String()) //这个发送者是 P2PKH 的，在BTC系统中虽然 P2WPKH 更推荐但是 P2PKH 也是很常见的

	netParams := chaincfg.TestNet3Params

	preMap := NewUtxoFromOutMap(map[wire.OutPoint]*UtxoSenderAmountTuple{
		*MustNewOutPoint("a06b4450eb63c8a245e799800e7554ef7a25d415874a957720bee79be8fc15e2", 1): NewUtxoSenderAmountTuple(NewAddressTuple("tb1q92kpf4hlj5khmdalshlz6602lvs8vcakxz8hzq"), 49560582),
	})

	param, err := NewCustomParamFromMsgTx(msgTx, preMap)
	require.NoError(t, err)
	require.NotNil(t, param)
	require.Equal(t, btcutil.Amount(580144), param.GetFee())

	require.NoError(t, param.VerifyMsgTxSign(msgTx, &netParams))
	t.Log("success")
	require.NoError(t, param.CheckMsgTxParam(msgTx, &netParams))
	t.Log("success")
}

func TestCustomParam_CheckMsgTxParam(t *testing.T) {
	const txHex = "0100000002b881daed6b9813086857122b1fcd1f89105cde69e779a61c1d29db241f6015cb000000006b483045022100c32d939e11ed0afd59d07ebf1580efc79eed0fbd1418fc637a30b9d3a2961e3702202f7b52da83bd65439058c0130de7458e47b7bfa47d9ab54a235c21a19c9d5044012102965a2272a740ca07daba5a20c9e27917f3ca491a673747f904b27e54e213ad5dfdffffff99b097802bf2357e487140bfdcd0c1206c1d50edcaaddf32fbfec6a339606790020000006b48304502210095b2f68b1f8dff1d61e3aa889c4a2169562603bd15ce00d51e0829e896f24ff702207aa2d328591189cc2b63f5784824c54a83d9a2e20580eaec3a03be67d001678b012102965a2272a740ca07daba5a20c9e27917f3ca491a673747f904b27e54e213ad5dfdffffff0200e1f505000000001976a914ea25a42febc13b31d72c2fe8c3b938de4979413b88ac36b9bb12000000001976a9143419c851dfd17e2645a613c4dcbe8bab58f20ea088ac00000000"

	msgTx, err := NewMsgTxFromHex(txHex)
	require.NoError(t, err)
	t.Log(msgTx.TxHash().String()) //这个发送者是 P2PKH 的，在BTC系统中虽然 P2WPKH 更推荐但是 P2PKH 也是很常见的

	netParams := dogecoin.MainNetParams

	preMap := NewUtxoFromOutMap(map[wire.OutPoint]*UtxoSenderAmountTuple{
		*MustNewOutPoint("cb15601f24db291d1ca679e769de5c10891fcd1f2b1257680813986bedda81b8", 0): NewUtxoSenderAmountTuple(NewAddressTuple("D9taZdfvonxSn8USudmhqhwvE7wt3aPW79"), 9230995),
		*MustNewOutPoint("90676039a3c6fefb32dfadcaed501d6c20c1d0dcbf4071487e35f22b8097b099", 2): NewUtxoSenderAmountTuple(NewAddressTuple("D9taZdfvonxSn8USudmhqhwvE7wt3aPW79"), 437264864),
	})

	param, err := NewCustomParamFromMsgTx(msgTx, preMap)
	require.NoError(t, err)
	require.NotNil(t, param)
	require.Equal(t, btcutil.Amount(32203325), param.GetFee())

	require.NoError(t, param.VerifyMsgTxSign(msgTx, &netParams))
	t.Log("success")
	require.NoError(t, param.CheckMsgTxParam(msgTx, &netParams))
	t.Log("success")
}

func TestCustomParam_CheckMsgTxParam_BTC(t *testing.T) {
	const txHex = "0100000000010348e6656b65e9b62bc91825046dd5ccb3480c45ba13194225bc4b8cf9525a9fd10100000000fdffffffc088aadf6e27b1665ba5b8d19421149404eb2c4614e2f7f1d0670574434a82d80000000000fdffffffcb4158811f3f3ec91d5eee5796d967deabcc6ebb8951a96a4ddf8edfb79f99f90000000000fdffffff02d399500c0000000022512069f7a5ba7d82d16357c5e30afa91e792a46a576b85f476e4ed9d5699bfab4345a593721d00000000160014873d31ab3269eac8cac60d0b5870928f324366fc02483045022100f3f1e514fcc30596fa9f608eae3aadbf41f1fb5f6eca88026181dd9de5972edd022051720c9bc3b89436b6c37d7acdd4b0143d285e83455a1e43d08119f982b23034012102c520fff259352d6ef8f7081b1530130b870afc53bbd542ff919a8a58eed3655c0248304502210098d1ad1b13ed0587dc32ca1006e2314665cd8940e8b552c57edb8cdd43b5ad9b0220513f1e5ec6aa33c241d05e2281c9e5d48984d53d028ca2f9f8bfca5173d628cd01210236812da64618b29e227b4e61d9e8b093ba28c7bec76a326d2942955df2796f8102483045022100cb9f555ab1e75e5bf2b3b9a1465f8bba2ba2da89957fafd693799878b4649d08022011e5e5065682d5d0acd1ef2925697e31a603862566bc44dab1f4a912e978f70001210236812da64618b29e227b4e61d9e8b093ba28c7bec76a326d2942955df2796f8100000000"

	msgTx, err := NewMsgTxFromHex(txHex)
	require.NoError(t, err)
	t.Log(msgTx.TxHash().String()) //这个发送者是 P2PKH 的，在BTC系统中虽然 P2WPKH 更推荐但是 P2PKH 也是很常见的

	netParams := chaincfg.MainNetParams

	preMap := NewUtxoFromOutMap(map[wire.OutPoint]*UtxoSenderAmountTuple{
		*MustNewOutPoint("d19f5a52f98c4bbc25421913ba450c48b3ccd56d042518c92bb6e9656b65e648", 1): NewUtxoSenderAmountTuple(NewAddressTuple("bc1qvuhjmgfr4kxye8eh63qvkv3yst950u8mye9fxh"), 700623171),
		*MustNewOutPoint("d8824a43740567d0f1f7e214462ceb0494142194d1b8a55b66b1276edfaa88c0", 0): NewUtxoSenderAmountTuple(NewAddressTuple("bc1q963tmm9tv9884k60puxc3syyld0xzte3duy9uc"), 17232),
		*MustNewOutPoint("f9999fb7df8edf4d6aa95189bb6eccabde67d99657ee5e1dc93e3f1f815841cb", 0): NewUtxoSenderAmountTuple(NewAddressTuple("bc1q963tmm9tv9884k60puxc3syyld0xzte3duy9uc"), 17191),
	})

	param, err := NewCustomParamFromMsgTx(msgTx, preMap)
	require.NoError(t, err)
	require.NotNil(t, param)
	require.Equal(t, btcutil.Amount(578), param.GetFee()) //因为 BTC 正式币很贵所以这里消耗的聪数比测试链更少些

	require.NoError(t, param.VerifyMsgTxSign(msgTx, &netParams))
	t.Log("success")
	require.NoError(t, param.CheckMsgTxParam(msgTx, &netParams))
	t.Log("success")
}

func TestCustomParam_CheckMsgTxParam_BTC_TXN(t *testing.T) {
	const txHex = "01000000000103898c193f03e1280c06a3ec68f9cc20d6c5a5a26b86585458b86f3a22ce8aee560600000000fdffffff8d7d18ec848627700ee0a679527f29bd9e1093999ef1a8a29e91b5288ddc0b840300000000fdffffff90f5367b24c83ba077961b5969a45d75a5a5928a3ddf3f66207420f2e4067cc50200000000fdffffff05c90101000000000016001478e2c75ddcf0f203ac56f82da9095f59a85fcadb80841e00000000001600143a922337fe73ba0efa861b3fd95bd0e59f78291740476600000000002200208a6e4b0bde302b3ae49d1478eea606e38b9e1d7b8b9f0d77ff15b0738cc9170c5f9e2a0100000000160014046ab7f047f3928f25d021d8d293baafe68d9dc6ffa31902000000001976a9149ce3e6737260f3d8648dec7d0fd0e37341079a3e88ac0400483045022100a97e3067907d9b252a2e0aa00836494a8099dea833dcc7de888d8446c8115f110220762e94b626464140c3286f1bad386cf39b0c543c58516650dc0f0ee0e4ac618f01483045022100efba529079aec73c4071d7772c00530e78c8eff237703f9392a89d8c07bbaef802203992e455d8e9f70f1575075cb1ac840407fd26b8bb7d6b8416ce244e13a1abd80169522103a1588b7607a8f65a42cddd0ee570bd83279d5011acdd0f65d676a78e2b7c4c662102d4cda5c294433b62c3c7bc4a93aa627b8b301f57fb98a35e1cdd89dc21b9e5f2210296be57dd5095bfbfb3139215f26eec111f742302514f31751d0b43ce694b662553ae0400483045022100bb62ef46767175f8f2126fda25093a558e8442c12cca513dcbd20cb490071d9202200a23b43e90b7dbbfc77003481652de0ed77e714a80af7ad3e5de0a9551ae95550147304402203e98e79da2a8571364cf9692ae7a479e5c846609275224b050419c4f39f6082202205544f06635b20c3be8628234dfb6b5f15b0d75eb068fcc79648411cf44c54fe7016952210296ed4fa0b2f551f5dad4df2efd591e7c0d3d065d8d79093630a749be71a96cd7210318474f5db5107b43dcbca16b9270c7005e7f1c463883fcc7f1d2acfa17f295772102f312f61107564868365b4c0af21587cf010e36b05eac080af5ffa70ae7581ec953ae04004830450221009604ff5eaf29ae4f59a4a4b9e9b5a284515c5cfb2d2a709344269125464cdbe202206909d50834a9f8cde36026affa1a81a8b2e412c2702914268dda3dac59e9cf500147304402207517b398b3a83623c868079caac84132ab6300bfaf95d2649a56e7b8f96032c3022036f19813aaa5790ee8181160d65d834a419754589f2a8afe2ce075c265ef3a6a01695221025590abd3c9ada7c64b03bf6241428d8fa78a1549a4afb94e64b3562f05ade2032103c70325aae80186e2af388206976066b0042916bc689c42b4e2d133107d196d522103374495ccf28bdbbcc165fb7423ad9b092af248b10f485ad3cf6401c1f71bd32d53ae00000000"

	msgTx, err := NewMsgTxFromHex(txHex)
	require.NoError(t, err)
	t.Log(msgTx.TxHash().String()) //这个发送者是 P2PKH 的，在BTC系统中虽然 P2WPKH 更推荐但是 P2PKH 也是很常见的

	netParams := chaincfg.MainNetParams

	preMap := NewUtxoFromOutMap(map[wire.OutPoint]*UtxoSenderAmountTuple{
		*MustNewOutPoint("56ee8ace223a6fb8585458866ba2a5c5d620ccf968eca3060c28e1033f198c89", 6): NewUtxoSenderAmountTuple(NewAddressTuple("bc1qrhut5t4g2wa2lf9h48fcth869h48khe0avxlq60vs5m0y2s8memq6fjt7r"), 31759346),
		*MustNewOutPoint("840bdc8d28b5919ea2a8f19e9993109ebd297f5279a6e00e70278684ec187d8d", 3): NewUtxoSenderAmountTuple(NewAddressTuple("bc1qlhqvxmpcqqw64r82h7sm89hn5n8g9p7w6y2m7cxkvtytavgthups57rnws"), 16616320),
		*MustNewOutPoint("c57c06e4f2207420663fdf3d8a92a5a5755da469591b9677a03bc8247b36f590", 2): NewUtxoSenderAmountTuple(NewAddressTuple("bc1q673p5npwsz78j4vdzqltsa4fvtvpteynudu446n56xrpy38tgrkqjmzwpk"), 15201401),
	})

	param, err := NewCustomParamFromMsgTx(msgTx, preMap)
	require.NoError(t, err)
	require.NotNil(t, param)
	require.Equal(t, btcutil.Amount(3076), param.GetFee()) //因为 BTC 正式币很贵所以这里消耗的聪数比测试链更少些

	require.NoError(t, param.VerifyMsgTxSign(msgTx, &netParams))
	t.Log("success")
	require.NoError(t, param.CheckMsgTxParam(msgTx, &netParams))
	t.Log("success")
}
