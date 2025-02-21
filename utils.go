package gobtcsign

import (
	"bytes"
	"encoding/hex"
	"math"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// GetTxHash returns the transaction hash from the signed transaction message.
// GetTxHash 通过签名后的交易信息获得交易哈希
func GetTxHash(msgTx *wire.MsgTx) string {
	return msgTx.TxHash().String()
}

// CvtMsgTxToHex converts the Go message body to a hex string used on the BTC chain.
// CvtMsgTxToHex 把go语言消息体转换为btc链上通用的hex字符串
func CvtMsgTxToHex(msgTx *wire.MsgTx) (string, error) {
	outTo := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	if err := msgTx.Serialize(outTo); err != nil {
		return "", errors.WithMessage(err, "wrong serialize")
	}
	txHex := hex.EncodeToString(outTo.Bytes())
	return txHex, nil
}

// NewMsgTxFromHex deserializes a transaction message body from a hex string.
// NewMsgTxFromHex deserializes a transaction message body from a hex string.
func NewMsgTxFromHex(txHex string) (*wire.MsgTx, error) {
	data, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, errors.WithMessage(err, "wrong decode data")
	}
	var msgTx = &wire.MsgTx{}
	err = msgTx.Deserialize(bytes.NewReader(data))
	if err != nil {
		return nil, errors.WithMessage(err, "wrong deserialize")
	}
	return msgTx, nil
}

// GetAddressPkScript generates the corresponding public key script (PkScript) from the address string.
// GetAddressPkScript 根据地址字符串生成对应的公钥脚本（PkScript），地址和公钥脚本是一对一的
func GetAddressPkScript(addressString string, netParams *chaincfg.Params) ([]byte, error) {
	address, err := btcutil.DecodeAddress(addressString, netParams)
	if err != nil {
		return nil, errors.WithMessage(err, "wrong decode-address")
	}
	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return nil, errors.WithMessage(err, "wrong get-pk-script")
	}
	return pkScript, nil
}

// MustNewAddress decodes the address string and panics if there is an error.
// MustNewAddress 根据地址字符串生成地址对象，如果出错则抛出异常
func MustNewAddress(addressString string, netParams *chaincfg.Params) btcutil.Address {
	address, err := btcutil.DecodeAddress(addressString, netParams)
	if err != nil {
		panic(errors.WithMessage(err, "wrong decode-address"))
	}
	return address
}

// MustGetPkScript generates the public key script (PkScript) from the address and panics if there is an error.
// MustGetPkScript 根据地址生成公钥脚本（PkScript），如果出错则抛出异常
func MustGetPkScript(address btcutil.Address) []byte {
	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		panic(errors.WithMessage(err, "wrong get-pk-script"))
	}
	return pkScript
}

// MustNewOutPoint creates a new OutPoint from the transaction hash and UTXO index, and panics if there is an error.
// MustNewOutPoint 根据交易哈希和UTXO索引创建新的OutPoint，如果出错则抛出异常
func MustNewOutPoint(srcTxHash string, utxoIndex uint32) *wire.OutPoint {
	//which tx the utxo from.
	utxoHash, err := chainhash.NewHashFromStr(srcTxHash)
	if err != nil {
		panic(errors.WithMessagef(err, "wrong param utxo-from-tx-hash=%s", srcTxHash))
	}
	return wire.NewOutPoint(
		utxoHash,  // 这个是收到 utxo 的交易哈希，即 utxo 是从哪里来的，配合位置索引序号构成唯一索引，就能确定是花的哪个utxo
		utxoIndex, // 这个是收到 utxo 的输出位置，比如一个交易中有多个输出，这里要选择输出的位置
	)
}

// NewInputOuts converts pkScripts and amounts to []*wire.TxOut.
// NewInputOuts 因为 SignParam 的成员里有 []*wire.TxOut 类型的前置输出字段，但教程常用的是 pkScripts [][]byte 和 amounts []int64 两个属性，因此这里写个转换逻辑
func NewInputOuts(pkScripts [][]byte, amounts []int64) []*wire.TxOut {
	size := max(len(pkScripts), len(amounts)) // must same size. so use the max size
	outs := make([]*wire.TxOut, 0, size)
	for idx := 0; idx < size; idx++ {
		outs = append(outs, wire.NewTxOut(amounts[idx], pkScripts[idx]))
	}
	return outs
}

// NewInputOutsV2 converts pkScripts and amounts to []*wire.TxOut using btcutil.Amount.
// NewInputOutsV2 因为 SignParam 的成员里有 []*wire.TxOut 类型的前置输出字段，但教程常用的是 pkScripts [][]byte 和 amounts []btcutil.Amount 两个属性，因此这里写个转换逻辑
func NewInputOutsV2(pkScripts [][]byte, amounts []btcutil.Amount) []*wire.TxOut {
	size := max(len(pkScripts), len(amounts)) // must same size. so use the max size
	outs := make([]*wire.TxOut, 0, size)
	for idx := 0; idx < size; idx++ {
		outs = append(outs, wire.NewTxOut(int64(amounts[idx]), pkScripts[idx]))
	}
	return outs
}

// GetMsgTxVSize returns the virtual size (v-size) of the signed transaction, which matches the value on the chain.
// GetMsgTxVSize 获得【签名后的】交易的大小，结果是 v-size 的，而且和链上的值相同
func GetMsgTxVSize(msgTx *wire.MsgTx) int {
	return int(math.Ceil(float64(3*msgTx.SerializeSizeStripped()+msgTx.SerializeSize()) / 4))
}

// GetRawTransaction retrieves the raw transaction from the client using the transaction hash.
// GetRawTransaction 通过交易哈希从客户端获取原始交易
func GetRawTransaction(client *rpcclient.Client, txHash string) (*btcjson.TxRawResult, error) {
	oneHash, err := chainhash.NewHashFromStr(txHash)
	if err != nil {
		return nil, errors.WithMessage(err, "wrong param-tx-hash")
	}
	return client.GetRawTransactionVerbose(oneHash)
}
