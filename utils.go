package gobtcsign

import (
	"bytes"
	"encoding/hex"
	"math"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/pkg/errors"
)

// GetTxHash 通过签名后的交易信息获得交易哈希
// 这个函数非常重要能够让你在发交易前就能知道哈希，这样就能在交易发出前就把信息存在数据库里，当交易发出后就可以去链上查找这笔交易
// 避免出现发了交易以后找不到的问题
// 理论上所有的链都能在交易发出以前就得到交易哈希，以方便程序员写出逻辑严密的代码，比如高并发和高可用情况下的收发交易
// 否则就很不好用，因此在做其它链时要多找找教程，在自己设计链时也要优先考虑提供这个功能
func GetTxHash(msgTx *wire.MsgTx) string {
	return msgTx.TxHash().String()
}

// CvtMsgTxToHex 把go语言消息体转换为btc链上通用的hex字符串
func CvtMsgTxToHex(msgTx *wire.MsgTx) (string, error) {
	outTo := bytes.NewBuffer(make([]byte, 0, msgTx.SerializeSize()))
	if err := msgTx.Serialize(outTo); err != nil {
		return "", errors.WithMessage(err, "wrong serialize")
	}
	txHex := hex.EncodeToString(outTo.Bytes())
	return txHex, nil
}

// NewMsgTxFromHex 通过交易信息的hex字符串，再反序列化出交易消息体
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

// GetAddressPkScript 根据地址字符串生成对应的公钥脚本（PkScript），地址和公钥脚本是一对一的
// 这个函数很重要，因为某些（少数的）函数的参数需要地址信息，而某些（多数的）函数需要公钥脚本信息
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

func MustNewAddress(addressString string, netParams *chaincfg.Params) btcutil.Address {
	address, err := btcutil.DecodeAddress(addressString, netParams)
	if err != nil {
		panic(errors.WithMessage(err, "wrong decode-address"))
	}
	return address
}

func MustGetPkScript(address btcutil.Address) []byte {
	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		panic(errors.WithMessage(err, "wrong get-pk-script"))
	}
	return pkScript
}

// IsDustOutputBtc 检查是不是灰尘输出，链不允许灰尘输出，避免到处都是粉尘
func IsDustOutputBtc(netParam *chaincfg.Params, address string, amount int64, dustLimitCoin float64, feePerKb btcutil.Amount) (bool, error) {
	dustLimitAmount, err := btcutil.NewAmount(dustLimitCoin)
	if err != nil {
		return false, errors.WithMessage(err, "wrong dust limit coin to amount")
	}
	if btcutil.Amount(amount) < dustLimitAmount { //当小于硬性灰尘数量时，就肯定是灰尘的
		return true, nil
	}
	pkScript, err := GetAddressPkScript(address, netParam)
	if err != nil {
		return false, errors.WithMessage(err, "wrong address->pk-script")
	}
	output := wire.NewTxOut(amount, pkScript)
	isDust := txrules.IsDustOutput(output, feePerKb)
	return isDust, nil
}

// IsDustOutputDoge 根据 https://github.com/dogecoin/dogecoin/blob/master/doc/fee-recommendation.md 这个看到只有两个规则
func IsDustOutputDoge(amount int64, dustLimitCoin float64) (bool, error) {
	dustLimitAmount, err := btcutil.NewAmount(dustLimitCoin)
	if err != nil {
		return false, errors.WithMessage(err, "wrong dust limit coin to amount")
	}
	if btcutil.Amount(amount) < dustLimitAmount { //当小于硬性灰尘数量时，就肯定是灰尘的
		return true, nil
	}
	return false, nil
}

// NewInputOuts 因为 SignParam 的成员里有 []*wire.TxOut 类型的前置输出字段
// 但教程常用的是 pkScripts [][]byte 和 amounts []int64 两个属性
// 因此这里写个转换逻辑
func NewInputOuts(pkScripts [][]byte, amounts []int64) []*wire.TxOut {
	size := max(len(pkScripts), len(amounts)) // must same size. so use the max size
	outs := make([]*wire.TxOut, 0, size)
	for idx := 0; idx < size; idx++ {
		outs = append(outs, wire.NewTxOut(amounts[idx], pkScripts[idx]))
	}
	return outs
}

func NewInputOutsV2(pkScripts [][]byte, amounts []btcutil.Amount) []*wire.TxOut {
	size := max(len(pkScripts), len(amounts)) // must same size. so use the max size
	outs := make([]*wire.TxOut, 0, size)
	for idx := 0; idx < size; idx++ {
		outs = append(outs, wire.NewTxOut(int64(amounts[idx]), pkScripts[idx]))
	}
	return outs
}

// GetMsgTxVSize 获得【签名后的】交易的大小，结果是 v-size 的，而且和链上的值相同
func GetMsgTxVSize(msgTx *wire.MsgTx) int {
	return int(math.Ceil(float64(3*msgTx.SerializeSizeStripped()+msgTx.SerializeSize()) / 4))
}
