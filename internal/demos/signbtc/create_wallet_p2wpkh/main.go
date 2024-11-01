package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func main() {
	netParams := &chaincfg.MainNetParams

	// 创建一个新的随机私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		log.Fatalf("随机私钥出错: %v", err)
	}
	// WIF（Wallet Import Format）私钥编码格式的类型
	privateWif, err := btcutil.NewWIF(privateKey, netParams, true)
	if err != nil {
		log.Fatalf("创建钱包引用格式出错: %v", err)
	}

	// 直接从私钥生成公钥
	pubKey := privateWif.PrivKey.PubKey()
	// 计算公钥哈希（P2WPKH使用的公钥哈希是公钥的SHA256和RIPEMD160哈希值）
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())
	// 创建P2WPKH地址
	witnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		log.Fatalf("创建P2WPKH地址出错: %v", err)
	}

	fmt.Println("私钥(WIF):", privateWif.String())
	fmt.Println("私钥(Hex):", hex.EncodeToString(privateKey.Serialize()))
	fmt.Println("P2WPKH地址:", witnessPubKeyHash.EncodeAddress())
	fmt.Println("地址网络名称:", netParams.Name)
}
