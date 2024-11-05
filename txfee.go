package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txrules"
	"github.com/btcsuite/btcwallet/wallet/txsizes"
	"github.com/pkg/errors"
	"github.com/yyle88/gobtcsign/internal/dusts"
)

type DustFee = dusts.DustFee

func NewDustFee() DustFee {
	return dusts.NewDustFee() //比特币没有软灰尘收费，这里配置个空的（因为doge里有，这里为了逻辑相通，而给个空的）
}

// EstimateTxFee 通过未签名且未找零的交易，预估出需要的费用
// 代码基本是仿照这里的 github.com/btcsuite/btcwallet/wallet/txauthor@v1.3.4/author.go 里面 NewUnsignedTransaction 的逻辑
// 由于是计算手续费的，因为这个交易里不应该包含找零的 output 信息，否则结果是无意义的
func EstimateTxFee(param *CustomParam, netParams *chaincfg.Params, change *ChangeTo, feeRatePerKb btcutil.Amount, dustFee DustFee) (btcutil.Amount, error) {
	maxSignedSize, err := EstimateTxSize(param, netParams, change)
	if err != nil {
		return 0, errors.WithMessage(err, "wrong estimate-tx-size")
	}
	outputs, err := param.GetOutputs(netParams)
	if err != nil {
		return 0, errors.WithMessage(err, "wrong get-outputs")
	}
	//有的链比如 DOGE_COIN 有软灰尘的概念，软灰尘需要消耗更高的手续费，而且这个手续费是不能协商的，而是必须交的，就得在这里交灰尘费
	maxRequiredFee := txrules.FeeForSerializeSize(feeRatePerKb, maxSignedSize) + dustFee.SumExtraDustFee(outputs)
	//但是请注意，input-output-maxFee 的结果还可能是个软灰尘，这时候就还得再增加找零的软灰尘费用，这个是后续逻辑需要考虑的
	return maxRequiredFee, nil
}

// EstimateTxSize 通过未签名的交易，预估出签名后交易体的大小，结果是 v-size 的，而且略微>=实际值
func EstimateTxSize(param *CustomParam, netParams *chaincfg.Params, change *ChangeTo) (int, error) {
	var scripts = make([][]byte, 0, len(param.VinList))
	for _, txIn := range param.VinList {
		pkScript, err := txIn.Sender.GetPkScript(netParams)
		if err != nil {
			return 0, errors.WithMessage(err, "wrong get-pk-script")
		}
		scripts = append(scripts, pkScript)
	}
	outputs, err := param.GetOutputs(netParams)
	if err != nil {
		return 0, errors.WithMessage(err, "wrong get-outputs")
	}
	return EstimateSize(scripts, outputs, change)
}

// EstimateSize 计算交易的预估大小（在最坏情况下的预估大小）
// 这个函数还是抄的 github.com/btcsuite/btcwallet/wallet/txauthor@v1.3.4/author.go 里面 NewUnsignedTransaction 的逻辑
// 详细细节见 https://github.com/btcsuite/btcwallet/blob/master/wallet/txauthor/author.go 这里的逻辑
// 是否填写找零信息，得依据 outputs 里面是否已经包含找零信息
func EstimateSize(scripts [][]byte, outputs []*wire.TxOut, change *ChangeTo) (int, error) {
	changeScriptSize, err := change.GetChangeScriptSize()
	if err != nil {
		return 0, errors.WithMessage(err, "wrong calculate-change-script-size")
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

type ChangeTo struct {
	PkScript []byte          //允许为空，当两者皆为空时表示没有找零输出
	AddressX btcutil.Address //允许为空，当两者皆为空时表示没有找零输出
}

func (T *ChangeTo) GetChangeScriptSize() (int, error) {
	if T.PkScript != nil { //优先使用找零脚本进行计算
		return CalculateChangePkScriptSize(T.PkScript)
	}
	if T.AddressX != nil { //其次使用找零地址进行计算
		return CalculateChangeAddressSize(T.AddressX)
	}
	return 0, nil //说明不需要找零输出，就返回0
}

func CalculateChangeAddressSize(address btcutil.Address) (int, error) {
	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return 0, errors.WithMessage(err, "wrong change_address")
	}
	return CalculateChangePkScriptSize(pkScript)
}

func CalculateChangePkScriptSize(pkScript []byte) (int, error) {
	var size int
	switch {
	case txscript.IsPayToPubKeyHash(pkScript):
		size = txsizes.P2PKHPkScriptSize
	case txscript.IsPayToScriptHash(pkScript):
		size = txsizes.NestedP2WPKHPkScriptSize
	case txscript.IsPayToWitnessPubKeyHash(pkScript), txscript.IsPayToWitnessScriptHash(pkScript):
		size = txsizes.P2WPKHPkScriptSize
	case txscript.IsPayToTaproot(pkScript):
		size = txsizes.P2TRPkScriptSize
	default:
		return 0, errors.New("UNSUPPORTED ADDRESS TYPE")
	}
	return size, nil
}
