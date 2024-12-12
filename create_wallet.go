package gobtcsign

import (
	"encoding/hex"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
)

// CreateWalletP2PKH generates a Bitcoin wallet using the P2PKH format.
// This function returns the wallet address and private key hex-string.
// CreateWalletP2PKH 使用 P2PKH 格式生成比特币钱包。
// 该函数返回钱包地址和私钥的十六进制格式。
func CreateWalletP2PKH(netParams *chaincfg.Params) (addressString string, privateKeyHex string, err error) {
	// Generate a new Bitcoin private key // 创建新的比特币私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", errors.WithMessage(err, "wrong to generate random private key")
	}

	// Generate the public key hash (SHA256 -> RIPEMD160) from the private key // 从私钥生成公钥哈希（SHA256 -> RIPEMD160）
	pubKeyHash := btcutil.Hash160(privateKey.PubKey().SerializeCompressed())

	// Create a Bitcoin address using the public key hash // 使用公钥哈希生成比特币地址
	addressPubKeyHash, err := btcutil.NewAddressPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		return "", "", errors.WithMessage(err, "wrong to create address from public key hash")
	}

	// Return the generated address and private key (hex-encoded) // 返回生成的地址和私钥（十六进制编码）
	addressString = addressPubKeyHash.EncodeAddress()
	privateKeyHex = hex.EncodeToString(privateKey.Serialize())
	return addressString, privateKeyHex, nil
}

// CreateWalletP2WPKH generates a Bitcoin wallet using the P2WPKH format.
// This function returns the wallet address and private key hex-string.
// CreateWalletP2WPKH 使用 P2WPKH 格式生成比特币钱包。
// 该函数返回钱包地址和私钥的十六进制格式。
func CreateWalletP2WPKH(netParams *chaincfg.Params) (addressString string, privateKeyHex string, err error) {
	// Generate a new Bitcoin private key // 创建新的比特币私钥
	privateKey, err := btcec.NewPrivateKey()
	if err != nil {
		return "", "", errors.WithMessage(err, "wrong to generate random private key")
	}

	// Encode the private key using Wallet Import Format (WIF) // 使用 WIF 格式编码私钥
	privateWif, err := btcutil.NewWIF(privateKey, netParams, true)
	if err != nil {
		return "", "", errors.WithMessage(err, "wrong to create Wallet Import Format (WIF) for private key")
	}

	// Get the public key // 获取公钥
	pubKey := privateWif.PrivKey.PubKey()

	// Compute the public key hash (SHA256 -> RIPEMD160) // 计算公钥哈希（SHA256 -> RIPEMD160）
	pubKeyHash := btcutil.Hash160(pubKey.SerializeCompressed())

	// Create a P2WPKH address using the public key hash // 使用公钥哈希生成 P2WPKH 地址
	witnessPubKeyHash, err := btcutil.NewAddressWitnessPubKeyHash(pubKeyHash, netParams)
	if err != nil {
		return "", "", errors.WithMessage(err, "wrong to create P2WPKH address")
	}

	// Return the generated address and private key (hex-encoded) // 返回生成的地址和私钥（十六进制编码）
	addressString = witnessPubKeyHash.EncodeAddress()
	privateKeyHex = hex.EncodeToString(privateKey.Serialize())
	return addressString, privateKeyHex, nil
}
