package gobtcsign

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
)

// CreateWalletP2PKH 随机创建个比特币钱包
// 你需要知道的是目前比特币的地址格式有5种，而这5种分别代表5个版本的，他们的签名逻辑各不相同
// 随着系统的升级以后可能还会有更多的版本出现
// 这里只是选择 P2PKH 这一种
// 其次是比特币分为正式网络和测试网络，他们的地址也不是互通的
// 需要在概念上区分清楚
func CreateWalletP2PKH(params *chaincfg.Params) (address string, private string, err error) {
	// 随机生成一个新的比特币私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", errors.WithMessage(err, "随机私钥失败")
	}
	// 通过私钥生成比特币地址的公钥
	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	// 通过公钥得到地址
	addressPubKeyHash, err := btcutil.NewAddressPubKeyHash(pubKeyHash, params)
	if err != nil {
		return "", "", errors.WithMessage(err, "创建地址失败")
	}
	addressString := addressPubKeyHash.EncodeAddress() // 转换为浏览器里常用的字符串的结果

	// 把私钥序列化得到二进制串
	priBytes := privateKey.Serialize()
	// 将私钥字节转换为 Hex 编码的字符串
	privateKeyHex := hex.EncodeToString(priBytes) // 转换为hex字符串
	return addressString, privateKeyHex, nil
}
