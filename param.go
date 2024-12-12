package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// BitcoinTxParams 这是客户自定义的参数类型，表示要转入和转出的信息
type BitcoinTxParams struct {
	VinList []VinType //要转入进BTC节点的
	OutList []OutType //要从BTC节点转出的-这里面通常包含1个目标（转账）和1个自己（找零）
	RBFInfo RBFConfig //详见RBF机制，通常是需要启用RBF以免交易长期被卡的
}

type VinType struct {
	OutPoint wire.OutPoint //UTXO的主要信息
	Sender   AddressTuple  //发送者信息，钱包地址或者公钥文本，二选一填写即可
	Amount   int64         //发送数量，因为这里不是浮点数，因此很明显这里传的是聪的数量
	RBFInfo  RBFConfig     //还是RBF机制，前面的是控制整个交易的，这里控制单个UTXO的
}

type OutType struct {
	Target AddressTuple //接收者信息，钱包地址和公钥文本，二选一填写即可
	Amount int64        //聪的数量
}

// CreateTxSignParams 根据用户的输入信息拼接交易
func (param *BitcoinTxParams) CreateTxSignParams(netParams *chaincfg.Params) (*SignParam, error) {
	var msgTx = wire.NewMsgTx(wire.TxVersion)

	//这是发送者和发送数量的列表，很明显，这是需要签名的关键信息，现在只把待签名信息收集起来
	var inputOuts = make([]*wire.TxOut, 0, len(param.VinList))
	for _, input := range param.VinList {
		pkScript, err := input.Sender.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong sender.address->pk-script")
		}
		inputOuts = append(inputOuts, wire.NewTxOut(input.Amount, pkScript))
	}

	//设置 vin 列表，当然这里拼装交易和签名是分离的，因此这里设置的是未签名的 utxo 信息。注意，这里需要跟前面的待签名信息位置序号相同
	for _, input := range param.VinList {
		utxo := input.OutPoint
		txIn := wire.NewTxIn(wire.NewOutPoint(&utxo.Hash, uint32(utxo.Index)), nil, nil)
		if txIn.Sequence != wire.MaxTxInSequenceNum { //这里做个断言，因为我后面的逻辑都是基于默认值是它而写的，假如默认值不是它就闹乌龙啦
			return nil, errors.Errorf("wrong tx_in.sequence default value: %v", txIn.Sequence)
		}
		// 查看是否需要启用 RBF 机制
		if seqNo := param.GetTxInputSequence(input); seqNo != wire.MaxTxInSequenceNum {
			txIn.Sequence = seqNo
		}
		msgTx.AddTxIn(txIn)
	}

	//设置 vout 列表，这个不需要签名，因此只要把目标地址和数量设置上就行
	for _, output := range param.OutList {
		pkScript, err := output.Target.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong target.address->pk-script")
		}
		msgTx.AddTxOut(wire.NewTxOut(output.Amount, pkScript))
	}
	return &SignParam{
		MsgTx:     msgTx,
		InputOuts: inputOuts, //这里它和 vin 的数量完全相同，而且位置序号也相同，最终签名时也需要确保位置相同
		NetParams: netParams,
	}, nil
}

func (param *BitcoinTxParams) GetOutputs(netParams *chaincfg.Params) ([]*wire.TxOut, error) {
	outputs := make([]*wire.TxOut, 0, len(param.OutList))
	for _, output := range param.OutList {
		pkScript, err := output.Target.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong target.address->pk-script")
		}
		outputs = append(outputs, wire.NewTxOut(output.Amount, pkScript))
	}
	return outputs, nil
}

func (param *BitcoinTxParams) GetTxInputSequence(input VinType) uint32 {
	// 当你确实是需要对每个交易单独设置RBF时，就可以在这里设置，单独设置到这个 vin 里面
	if seqNo := input.RBFInfo.GetSequence(); seqNo != wire.MaxTxInSequenceNum { //启用RBF机制，精确的RBF逻辑
		return seqNo
	}
	// 这里不设置也行，设置是为了启用 RBF 机制，设置到全部 vin 里面，当然前面的 RBF 会优先设置
	if seqNo := param.RBFInfo.GetSequence(); seqNo != wire.MaxTxInSequenceNum { //启用RBF机制，粗放的RBF逻辑
		// RBF (Replace-By-Fee) 是比特币网络中的一种机制。搜索官方的 “RBF” 即可得到你想要的知识
		// 简单来说 RBF 就是允许使用相同 utxo 发两次不同的交易，但只有其中的一笔能生效
		// 在启用 RBF 时发第二笔交易会报错，而允许重发时，发第二笔以后这两笔交易都会成为待打包状态，哪笔会打包和确认得看链上的打包情况
		// 通常，序列号设置为较高的值（如0xfffffffd），表示交易是可替换的
		// 因此，推荐的设置就是 txIn.Sequence = wire.MaxTxInSequenceNum - 2
		// 当然，设置为 0，1，2，3 也是可以的，只不过看着不太专业，推荐还是前面的 `0xfffffffd` 序列号
		// 理论上每个 txIn 都有独立的序列号，但是在业务中通常就是某个交易里的所有 txIn 使用相同的序列号，这样便于写CRUD逻辑
		return seqNo
	}
	// 当都没有设置的时候，就使用默认值就行
	return wire.MaxTxInSequenceNum
}

// GetInputList 把拼交易的参数转换为验签的参数
func (param *BitcoinTxParams) GetInputList() []*VerifyTxInputParam {
	var inputList = make([]*VerifyTxInputParam, 0, len(param.VinList))
	for _, x := range param.VinList {
		inputList = append(inputList, &VerifyTxInputParam{
			Sender: AddressTuple{
				Address:  x.Sender.Address,
				PkScript: x.Sender.PkScript,
			},
			Amount: x.Amount,
		})
	}
	return inputList
}

// GetFee 全部输入和全部输出的差额，即交易的费用
func (param *BitcoinTxParams) GetFee() btcutil.Amount {
	var sum int64
	for _, v := range param.VinList {
		sum += v.Amount
	}
	for _, v := range param.OutList {
		sum -= v.Amount
	}
	return btcutil.Amount(sum)
}

// GetChangeAmountWithFee 根据交易费用计算出找零数量
func (param *BitcoinTxParams) GetChangeAmountWithFee(fee btcutil.Amount) btcutil.Amount {
	return param.GetFee() - fee
}

func (param *BitcoinTxParams) EstimateTxSize(netParams *chaincfg.Params, change *ChangeTo) (int, error) {
	return EstimateTxSize(param, netParams, change)
}

func (param *BitcoinTxParams) EstimateTxFee(netParams *chaincfg.Params, change *ChangeTo, feeRatePerKb btcutil.Amount, dustFee DustFee) (btcutil.Amount, error) {
	return EstimateTxFee(param, netParams, change, feeRatePerKb, dustFee)
}

// NewCustomParamFromMsgTx 这里提供简易的逻辑把交易的原始参数再拼回来
// 以校验参数和校验签名等信息
// 因此该函数的主要作用是校验
// 首先拿到已签名(待发送/已发送)的交易的 hex 数据，接着使用 NewMsgTxFromHex 即可得到交易数据
// 接着使用此函数再反拼出原始参数，检查交易的费用，接着再检查签名
// 第二个参数是设置如何获取前置输出的
// 通常是使用 客户端 请求获取前置输出，但也可以使用map把前置输出存起来，因此使用 interface 获取前置输出，提供两种实现方案
// 在项目中推荐使用 rpc 获取，这样就很方便，而在单元测试中则只需要通过 map 预先配置就行，避免网络请求也避免暴露节点配置
func NewCustomParamFromMsgTx(msgTx *wire.MsgTx, preImp GetUtxoFromInterface) (*BitcoinTxParams, error) {
	var vinList = make([]VinType, 0, len(msgTx.TxIn))
	for _, vin := range msgTx.TxIn {
		costUtxo := vin.PreviousOutPoint

		utxoFrom, err := preImp.GetUtxoFrom(costUtxo)
		if err != nil {
			return nil, errors.WithMessage(err, "get-utxo-from")
		}

		vinList = append(vinList, VinType{
			OutPoint: *wire.NewOutPoint(&costUtxo.Hash, costUtxo.Index),
			Sender:   *utxoFrom.sender,
			Amount:   utxoFrom.amount,
			RBFInfo:  *NewRBFConfig(vin.Sequence),
		})
	}

	var outList = make([]OutType, 0, len(msgTx.TxOut))
	for _, out := range msgTx.TxOut {
		outList = append(outList, OutType{
			Target: AddressTuple{PkScript: out.PkScript},
			Amount: out.Value,
		})
	}

	param := &BitcoinTxParams{
		VinList: vinList,
		OutList: outList,
		RBFInfo: *NewRBFNotUse(), //这里是不需要的，因为各个输入里将会有RBF的全部信息
	}
	return param, nil
}

// VerifyMsgTxSign 使用这个检查签名是否正确
func (param *BitcoinTxParams) VerifyMsgTxSign(msgTx *wire.MsgTx, netParams *chaincfg.Params) error {
	inputsItem, err := param.GetVerifyTxInputsItem(netParams)
	if err != nil {
		return errors.WithMessage(err, "wrong get-inputs")
	}
	if err := VerifySignV3(msgTx, inputsItem); err != nil {
		return errors.WithMessage(err, "wrong verify-sign")
	}
	return nil
}

func (param *BitcoinTxParams) GetVerifyTxInputsItem(netParams *chaincfg.Params) (*VerifyTxInputsType, error) {
	var res = &VerifyTxInputsType{
		PkScripts: make([][]byte, 0, len(param.VinList)),
		InAmounts: make([]btcutil.Amount, 0, len(param.VinList)),
	}
	for _, vin := range param.VinList {
		pkScript, err := vin.Sender.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong sender.address->pk-script")
		}
		res.PkScripts = append(res.PkScripts, pkScript)

		res.InAmounts = append(res.InAmounts, btcutil.Amount(vin.Amount))
	}
	return res, nil
}
