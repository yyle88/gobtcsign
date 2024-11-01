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
func CreateWalletP2PKH(netParams *chaincfg.Params) (address string, private string, err error) {
	// 随机生成一个新的比特币私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", errors.WithMessage(err, "随机私钥出错")
	}
	// 通过私钥生成比特币地址的公钥
	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())
	// 通过公钥得到地址
	addressPubKeyHash, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		return "", "", errors.WithMessage(err, "创建地址出错")
	}
	address = addressPubKeyHash.EncodeAddress() // 转换为浏览器里常用的字符串的结果
	private = hex.EncodeToString(privateKey.Serialize())
	return address, private, nil
}

func CreateWalletP2WPKH(netParams *chaincfg.Params) (address string, private string, err error) {
	// 创建一个新的随机私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", errors.WithMessage(err, "随机私钥出错")
	}
	// WIF（Wallet Import Format）私钥编码格式的类型
	privateWif, err := btcutil.NewWIF(privateKey, netParams, true)
	if err != nil {
		return "", "", errors.WithMessage(err, "创建钱包引用格式出错")
	}
	// 直接从私钥生成公钥
	pubKey := privateWif.PrivKey.PubKey()
	// 计算公钥哈希（P2WPKH使用的公钥哈希是公钥的SHA256和RIPEMD160哈希值）
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	// 创建P2WPKH地址
	witnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		return "", "", errors.WithMessage(err, "创建P2WPKH地址出错")
	}
	address = witnessPubKeyHash.EncodeAddress()
	private = hex.EncodeToString(privateKey.Serialize())
	return address, private, nil
}
