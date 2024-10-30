package gobtcsign

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/wire"
	"github.com/pkg/errors"
)

// CustomParam 这是客户自定义的参数类型，表示要转入和转出的信息
type CustomParam struct {
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

func MustNewOutPoint(srcTxHash string, utxoIndex uint32) *wire.OutPoint {
	//which tx the utxo from.
	utxoHash, err := chainhash.NewHashFromStr(srcTxHash)
	if err != nil {
		panic(errors.WithMessagef(err, "wrong param utxo-from-tx-hash=%s", srcTxHash))
	}
	return wire.NewOutPoint(
		utxoHash,  //这个是收到 utxo 的交易哈希，即 utxo 是从哪里来的，配合位置索引序号构成唯一索引，就能确定是花的哪个utxo
		utxoIndex, //这个是收到 utxo 的输出位置，比如一个交易中有多个输出，这里要选择输出的位置
	)
}

type OutType struct {
	Target AddressTuple //接收者信息，钱包地址和公钥文本，二选一填写即可
	Amount int64        //聪的数量
}

type AddressTuple struct {
	Address  string //钱包地址 和 公钥脚本 二选一填写即可
	PkScript []byte //公钥脚本 和 钱包地址 二选一填写即可 PkScript（Public Key Script）在拼装交易和签名时使用
}

func NewAddressTuple(address string) AddressTuple {
	return AddressTuple{
		Address:  address,
		PkScript: nil, //这里 address 和 pk-script 是二选一的，因此不设，在后续的逻辑里会根据地址获得 pk-script 信息
	}
}

// GetPkScript 获得公钥文本，当公钥文本存在时就用已有的，否则就根据地址计算
func (one *AddressTuple) GetPkScript(netParams *chaincfg.Params) ([]byte, error) {
	if len(one.PkScript) > 0 {
		return one.PkScript, nil
	}
	return GetAddressPkScript(one.Address, netParams)
}

type RBFConfig struct {
	AllowRBF bool   //当需要RBF时需要设置，推荐启用RBF发交易，否则，当手续费过低时交易会卡在节点的内存池里
	Sequence uint32 //这是后来出的功能，RBF，使用更高手续费重发交易用的，当BTC交易发出到系统以后假如没人打包（手续费过低时），就可以增加手续费覆盖旧的交易
}

func NewRBFActive() RBFConfig {
	return RBFConfig{
		AllowRBF: true,
		Sequence: wire.MaxTxInSequenceNum - 2, // recommended sequence BTC推荐的默认启用RBF的就是这个数
	}
}

func NewRBFNotUse() RBFConfig {
	return RBFConfig{
		AllowRBF: false, //当两个元素都为零值时表示不启用RBF机制
		Sequence: 0,     //当两个元素都为零值时表示不启用RBF机制
	}
}

func (cfg *RBFConfig) GetSequence() uint32 {
	if cfg.AllowRBF || cfg.Sequence > 0 { //启用RBF机制，精确的RBF逻辑
		return cfg.Sequence
	}
	return wire.MaxTxInSequenceNum //当两个元素都为零值时表示不启用RBF机制-因此这里使用默认的最大值表示不启用
}

// GetSignParam 根据用户的输入信息拼接交易
func (param *CustomParam) GetSignParam(netParams *chaincfg.Params) (*SignParam, error) {
	var msgTx = wire.NewMsgTx(wire.TxVersion)
	var pkScripts [][]byte
	var amounts []int64
	for _, input := range param.VinList {
		pkScript, err := input.Sender.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong sender.address->pk-script")
		}
		pkScripts = append(pkScripts, pkScript)
		amounts = append(amounts, input.Amount)

		utxo := input.OutPoint
		txIn := wire.NewTxIn(wire.NewOutPoint(&utxo.Hash, uint32(utxo.Index)), nil, nil)
		if txIn.Sequence != wire.MaxTxInSequenceNum { //这里做个断言，因为我后面的逻辑都是基于默认值是它而写的，假如默认值不是它就闹乌龙啦
			return nil, errors.Errorf("wrong tx_in.sequence default value: %v", txIn.Sequence)
		}
		// 查看是否需要启用 RBF 机制
		if seqNo := param.GetTxInSequenceNum(input); seqNo != wire.MaxTxInSequenceNum {
			txIn.Sequence = seqNo
		}
		msgTx.AddTxIn(txIn)
	}
	for _, output := range param.OutList {
		pkScript, err := output.Target.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong target.address->pk-script")
		}
		msgTx.AddTxOut(wire.NewTxOut(output.Amount, pkScript))
	}
	return &SignParam{
		MsgTx:     msgTx,
		PkScripts: pkScripts,
		Amounts:   amounts,
		NetParams: netParams,
	}, nil
}

func (param *CustomParam) GetTxInSequenceNum(input VinType) uint32 {
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
func (param *CustomParam) GetInputList() []*VerifyTxInputParam {
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
