package gosignbtc

import (
	"encoding/hex"
	"reflect"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// CustomParam 这是客户自定义的参数类型
type CustomParam struct {
	VinList []VinType
	OutList []OutType
}

type VinType struct {
	Utxo     *wire.OutPoint
	PkScript []byte
	Amount   int64
}

type OutType struct {
	Address string
	Amount  int64
}

// SignParam 这是系统内的参数类型，通过客户自定义的类型能得到待签名的信息
type SignParam struct {
	msgTx     *wire.MsgTx
	pkScripts [][]byte
	amounts   []int64
	netParams *chaincfg.Params
}

// GetMsgTx 通常用户是不需要修改这个值的，直接就能签名，但是假如用户需要修改的话也可以通过这个函数得到它
func (p *SignParam) GetMsgTx() *wire.MsgTx {
	return p.msgTx
}

// NewSignParam 根据用户的输入信息拼接交易
func NewSignParam(param CustomParam, netParams *chaincfg.Params) (*SignParam, error) {
	var msgTx = wire.NewMsgTx(wire.TxVersion)
	var pkScripts [][]byte
	var amounts []int64
	for _, input := range param.VinList {
		pkScripts = append(pkScripts, input.PkScript)
		amounts = append(amounts, input.Amount)

		utxo := input.Utxo
		txIn := wire.NewTxIn(wire.NewOutPoint(&utxo.Hash, uint32(utxo.Index)), nil, nil)
		//txIn.Sequence = param.UtxoTxInSequence //这里不设置也行，设置是为了重发交易
		msgTx.AddTxIn(txIn)
	}
	for _, output := range param.OutList {
		address, err := btcutil.DecodeAddress(output.Address, netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong encrypt.decode_address")
		}
		pkScript, err := txscript.PayToAddrScript(address)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong encrypt.pay_to_addr_script")
		}
		msgTx.AddTxOut(wire.NewTxOut(output.Amount, pkScript))
	}
	return &SignParam{
		msgTx:     msgTx,
		pkScripts: pkScripts,
		amounts:   amounts,
		netParams: netParams,
	}, nil
}

// Sign 根据钱包地址和钱包私钥签名
func Sign(fromAddress string, privateKeyHex string, param *SignParam) error {
	privKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return errors.WithMessage(err, "wrong decode private key string")
	}
	privKey, pubKey := btcec.PrivKeyFromBytes(privKeyBytes)

	walletAddress, err := btcutil.DecodeAddress(fromAddress, param.netParams)
	if err != nil {
		return errors.WithMessage(err, "wrong from_address")
	}

	switch address := walletAddress; address.(type) {
	case *btcutil.AddressPubKeyHash: //请参考 txscript.PubKeyHashTy 的签名逻辑
		//检查钱包的地址是不是压缩的，有压缩和不压缩两种格式的地址，都是可以用的
		compress, err := CheckPKHAddressIsCompress(param.netParams, pubKey, fromAddress)
		if err != nil {
			return errors.WithMessage(err, "wrong sign check_from_address_is_compress")
		}
		//根据是否压缩选择不同的签名逻辑
		if err := SignP2PKH(param, privKey, compress); err != nil {
			return errors.WithMessage(err, "wrong sign")
		}
	default: //其它钱包类型暂不支持
		return errors.Errorf("From地址 %s 属于 %s 类型, 类型错误", address, reflect.TypeOf(address).String()) //倒是没必要支持太多的类型
	}
	return nil
}

func CheckPKHAddressIsCompress(defaultNet *chaincfg.Params, publicKey *btcec.PublicKey, fromAddress string) (bool, error) {
	for _, isCompress := range []bool{true, false} {
		var pubKeyHash []byte
		if isCompress {
			pubKeyHash = btcutil.Hash160(publicKey.SerializeCompressed())
		} else {
			pubKeyHash = btcutil.Hash160(publicKey.SerializeUncompressed())
		}

		address, err := btcutil.NewAddressPubKeyHash(pubKeyHash, defaultNet)
		if err != nil {
			return isCompress, errors.Errorf("error=%v when is_compress=%v", err, isCompress)
		}
		if address.EncodeAddress() == fromAddress {
			return isCompress, nil
		}
	}
	return false, errors.Errorf("unknown address type. address=%s", fromAddress)
}

func SignP2PKH(signParam *SignParam, privKey *btcec.PrivateKey, compress bool) error {
	var (
		msgTx     = signParam.msgTx
		pkScripts = signParam.pkScripts
		amounts   = signParam.amounts
	)

	for idx := range msgTx.TxIn {
		// 使用私钥对交易输入进行签名
		// 在大多数情况下，使用压缩公钥是可以接受的，并且更常见。压缩公钥可以减小交易的大小，从而降低交易费用，并且在大多数情况下，与非压缩公钥相比，安全性没有明显的区别
		signatureScript, err := txscript.SignatureScript(msgTx, idx, pkScripts[idx], txscript.SigHashAll, privKey, compress)
		if err != nil {
			return errors.Errorf("wrong signature_script. index=%d error=%v", idx, err)
		}
		msgTx.TxIn[idx].SignatureScript = signatureScript
	}
	return VerifyP2PKHSign(msgTx, pkScripts, amounts)
}

func VerifyP2PKHSign(tx *wire.MsgTx, pkScripts [][]byte, amounts []int64) error {
	for idx := range tx.TxIn { // 这段代码的作用是创建和执行脚本引擎，用于验证指定的脚本是否有效。如果脚本验证失败，则返回错误信息。这在比特币交易的验证过程中非常重要，以确保交易的合法性和安全性。
		vm, err := txscript.NewEngine(pkScripts[idx], tx, idx, txscript.StandardVerifyFlags, nil, nil, amounts[idx], nil)
		if err != nil {
			return errors.Errorf("wrong vm. index=%d error=%v", idx, err)
		}
		if err = vm.Execute(); err != nil {
			return errors.Errorf("wrong vm execute. index=%d error=%v", idx, err)
		}
	}
	return nil
}
