package gobtcsign

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestVerifyTx_DOGE_testnet(t *testing.T) {
	const txHex = "0100000001d7c8f1b28cab162a517889c2d66f9194446cdc3a571bd7f8ffdbc15b16a0970a000000006b483045022100bba2e1fd90d763d77775f7b4488841846fa307c3b203355048694ae1f048027d0220191e32dc86ad41536b93791b6070bcedd2488dc2c8445572b018424cf9cb91c90121031c983815f246c81e22e824901143c407560047227394d8c39fc754d45fec2770fdffffff0280969800000000001976a914b92052765b0007a3e3f0375b34430d7e8df695ea88acf135f503000000001976a91496c179f468add293a9cfaebc127a62d09bcb57c288ac00000000"

	msgTx, err := NewMsgTxFromTxHex(txHex)
	require.NoError(t, err)
	t.Log(msgTx.TxHash().String())

	const address = "nhwHQ29uEKnHeiWNLpU2zVHtcrZjaLKaHF"
	netParams := &dogecoin.TestNetParams
	pkScript, err := GetAddressPkScript(address, netParams)
	require.NoError(t, err)

	require.NoError(t, VerifyP2PKHSignV2(msgTx, []*VerifyTxInputParam{
		&VerifyTxInputParam{
			Sender: AddressTuple{
				Address:  address,
				PkScript: pkScript,
			},
			Amount: 0, //P2PKH 签名不将 amount 包含在生成的签名哈希中，因此也不验证它，随便填都行
		},
	}, netParams))
}

func TestVerifyTx_DOGE_mainnet(t *testing.T) {
	const txHex = "01000000019cd6fd32cca16c3745ffe3d4a952ae7463df930f899f1b52424dffb4932d6fdb010000006a473044022013dbb0fe1678ae1f74fa3f861456f5b0396211b1866d202242d0c4becd48429e02202b0e55c020f8f2c05a40b192fd1c3ea20fe8f6299bbbf9bf87a41d14c62067200121023c3520f4835bb87b9580bbdc6cd364072b45167c774c4556c15c580cc5fcba7bffffffff02a0862d46010000001976a914df20bc04fd348aa2e7399d677de8a3eb742bb27488ac38e5f46bf30a00001976a9141303e3c789d3b835d95c8b6be590ad15af23c90b88ac00000000"

	msgTx, err := NewMsgTxFromTxHex(txHex)
	require.NoError(t, err)
	t.Log(msgTx.TxHash().String())

	const address = "D6se3Ajq9mF8YD4p7jSXwZxewT6ePsnea6"
	netParams := &dogecoin.MainNetParams
	pkScript, err := GetAddressPkScript(address, netParams)
	require.NoError(t, err)

	require.NoError(t, VerifyP2PKHSignV2(msgTx, []*VerifyTxInputParam{
		&VerifyTxInputParam{
			Sender: AddressTuple{
				Address:  address,
				PkScript: pkScript,
			},
			Amount: 0, //P2PKH 签名不将 amount 包含在生成的签名哈希中，因此也不验证它，随便填都行
		},
	}, netParams))
}

func TestVerifyTx_DOGE_mainnet_2(t *testing.T) {
	const txHex = "0100000003ecd4d9d0bfa4d2126d08d17686b069d0a36901b5650023d02c86dd0ad6b92264000000006a473044022072d8fb163976d8aebad8e413e0cc958022a0b481f5d4db3a5d111a5ef00cc63802200f25c83f9f44425799a9ca86e992a994ff2db8c0f0e80523c0dd4753fd13fff601210338a2c85b6cf704e6b2eab54c985b78f7aacde132082d4ff5c0d62d3cd08da6a4ffffffff8f6e68ae9eab24647ce54b9e275b2a8bef25f2865440b34358388720ad2ad073010000006b483045022100e8d2eec625ca70c98e2043cae9232420eed0cc2166e4bcebe53debfe8064306f02207740703a9cfeb7967dd2e8c62fb7957fe599a0bf91998623d96e85df50e0dfa0012102cc9c56c656f549e5f88dcf3a9d1bc54e6638cfe69b379f6d41a5d1f756376ca4ffffffff23da7350f40dfd635284552202fcb801289afa2d906a5af6ca33f1e02ffccf48060000006b483045022100cf5b201fd538c76f8a82ccbc19c04d975cee249fbe24c49e27fa0ad78e86a5c302200723c0e5b0f66d81b40aec6892462700a845d9490d40b2d91ff1f4a51c97312e012102e272a9019a284d1e1db8adada953b4bb2b2eebc9264e45b3886a5acdb6f7b049ffffffff0294d7ca41000000001976a914551137476a06cfdedd91279f60692682b3f3d68b88ac043b10ce7be60a001976a9140be4c25349ee33a3f3d9674fdd31618918cacd4588ac00000000"

	msgTx, err := NewMsgTxFromTxHex(txHex)
	require.NoError(t, err)
	t.Log(msgTx.TxHash().String())

	// 请看这里 https://sochain.com/tx/DOGE/aa9188f81b8114314318bd75fab275f40022ff1e2a83ea2c84b809ef2d562a29
	// 你需要做的就是把看到的 inputs 拷贝下来
	var addresses = []string{
		"DCYWbjW1rRyWnXSKnjp8tuQYUPsv9Yj1cT",
		"D8guyBsLjbLgnXQBsxtDgnbtwAh662aJGb",
		"DJnqNZJn21HKJjEzjLTTWToBpvZPjomTnR",
	}

	require.NoError(t, VerifyP2PKHSignV3(msgTx, addresses, &dogecoin.MainNetParams))
}

func TestVerifyTx_DOGE_mainnet_3(t *testing.T) {
	const txHex = "0100000002d2b3941020b7a48057bc8fbab13f45f681e1717ddadab6b83ac2dbd926cb914d000000006a473044022021a6ad19b23b1dc0be9914d8fd49af932f59e26f0a113ba225083199dc1afa9002206a6489794429bd86c90849d03b54230203500c45bd3532000d8c0c97cacbc433012103e2c5018631f5960a9c0c988d7b5355048bb8ff2a772c69a59167208f1a2a9896ffffffffe72683d7ed5c5cc9918823031131ad572e9b190c5ced1008a57838f952a77ce2010000006a47304402205c3447fda59028ff08f10809b6003af7d94ba84d52280d6362acaca8fbd752a8022060d40242de21b2dfd3a6ac0b6d0b37799b53721d69fb90c9e461c30ac9585221012103e2c5018631f5960a9c0c988d7b5355048bb8ff2a772c69a59167208f1a2a9896ffffffff0200e9587278e60a001976a91451368824f9f5badcf889c024331cc45d2087558a88ac80439bd0b56f00001976a914869d6c848672f73fcb6b357b5650e780bad2106988ac00000000"

	msgTx, err := NewMsgTxFromTxHex(txHex)
	require.NoError(t, err)
	t.Log(msgTx.TxHash().String())

	//请看这里 https://www.oklink.com/zh-hans/doge/tx/6422b9d60add862cd0230065b50169a3d069b08676d1086d12d2a4bfd0d9d4ec
	//你需要做的就是把 输入列表 拷贝下来
	var addresses = []string{
		"DHQsfy66JsYSnwjCABFN6NNqW4kHQe63oU",
		"DHQsfy66JsYSnwjCABFN6NNqW4kHQe63oU",
	}

	require.NoError(t, VerifyP2PKHSignV3(msgTx, addresses, &dogecoin.MainNetParams))
}
