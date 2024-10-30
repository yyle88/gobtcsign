package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/pkg/errors"
)

// CalculateMsgTxFee 通过未签名且未找零的交易，预估出需要的费用
// 代码基本是仿照这里的 github.com/btcsuite/btcwallet/wallet/txauthor@v1.3.4/author.go 里面 NewUnsignedTransaction 的逻辑
// 需要特别注意的是，这里只是个使用的样例
// 由于是计算手续费的，因为这个交易里不应该包含找零的 output 信息，否则结果是无意义的
func CalculateMsgTxFee(msgTx *wire.MsgTx, changeAddress btcutil.Address, feeRatePerKb btcutil.Amount, dustFeeConfig DustFee) (btcutil.Amount, error) {
	maxSignedSize, err := CalculateMsgTxSize(msgTx, changeAddress)
	if err != nil {
		return 0, errors.WithMessage(err, "calculate size is wrong")
	}
	//有的链比如 DOGE_COIN 有软灰尘的概念，软灰尘需要消耗更高的手续费，而且这个手续费是不能协商的，而是必须交的，就得在这里交灰尘费
	maxRequiredFee := txrules.FeeForSerializeSize(feeRatePerKb, maxSignedSize) + dustFeeConfig.SumExtraDustsFee(msgTx.TxOut)
	//但是请注意，input-output-maxFee 的结果还可能是个软灰尘，这时候就还得再增加找零的软灰尘费用，这个是后续逻辑需要考虑的
	return maxRequiredFee, nil
}

// CalculateMsgTxSize 计算交易的预估大小（在最坏情况下的预估大小）
// 这个函数还是抄的 github.com/btcsuite/btcwallet/wallet/txauthor@v1.3.4/author.go 里面 NewUnsignedTransaction 的逻辑
// 这个函数也是个样例，主要是演示如何预估大小，当然也是可以直接使用的
// 这个交易里不应该包含找零信息
func CalculateMsgTxSize(msgTx *wire.MsgTx, changeAddress btcutil.Address) (int, error) {
	var scripts = make([][]byte, 0, len(msgTx.TxIn))
	for _, txIn := range msgTx.TxIn {
		scripts = append(scripts, txIn.SignatureScript)
	}
	return CalculateSize(scripts, msgTx.TxOut, changeAddress)
}

// CalculateSize 计算交易的预估大小（在最坏情况下的预估大小）
// 这个函数的参数详见前面函数的调用
func CalculateSize(scripts [][]byte, outputs []*wire.TxOut, changeAddress btcutil.Address) (int, error) {
	changeScriptSize, err := CalculateChangeAddressSize(changeAddress)
	if err != nil {
		return 0, errors.WithMessage(err, "calculate change script size")
	}

	// We count the types of inputs, which we'll use to estimate
	// the vsize of the transaction.
	var nested, p2wpkh, p2tr, p2pkh int
	for _, pkScript := range scripts {
		switch {
		// If this is a p2sh output, we assume this is a
		// nested P2WKH.
		case txscript.IsPayToScriptHash(pkScript):
			nested++
		case txscript.IsPayToWitnessPubKeyHash(pkScript):
			p2wpkh++
		case txscript.IsPayToTaproot(pkScript):
			p2tr++
		default:
			p2pkh++
		}
	}

	// 仿照这个函数 txauthor.NewUnsignedTransaction() 里的预估逻辑
	maxSignedSize := txsizes.EstimateVirtualSize(
		p2pkh, p2tr, p2wpkh, nested, outputs, changeScriptSize,
	)
	return maxSignedSize, nil
}

func CalculateChangeAddressSize(address btcutil.Address) (int, error) {
	script, err := txscript.PayToAddrScript(address)
	if err != nil {
		return 0, errors.WithMessage(err, "wrong change_address")
	}
	var size int
	switch {
	case txscript.IsPayToPubKeyHash(script):
		size = txsizes.P2PKHPkScriptSize
	case txscript.IsPayToScriptHash(script):
		size = txsizes.NestedP2WPKHPkScriptSize
	case txscript.IsPayToWitnessPubKeyHash(script), txscript.IsPayToWitnessScriptHash(script):
		size = txsizes.P2WPKHPkScriptSize
	case txscript.IsPayToTaproot(script):
		size = txsizes.P2TRPkScriptSize
	default:
		return 0, errors.New("UNSUPPORTED ADDRESS TYPE")
	}
	return size, nil
}
