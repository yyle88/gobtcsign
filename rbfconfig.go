package gobtcsign

import "github.com/btcsuite/btcd/wire"

type RBFConfig struct {
	AllowRBF bool   //当需要RBF时需要设置，推荐启用RBF发交易，否则，当手续费过低时交易会卡在节点的内存池里
	Sequence uint32 //这是后来出的功能，RBF，使用更高手续费重发交易用的，当BTC交易发出到系统以后假如没人打包（手续费过低时），就可以增加手续费覆盖旧的交易
}

func NewRBFConfig(sequence uint32) *RBFConfig {
	return &RBFConfig{
		AllowRBF: sequence != wire.MaxTxInSequenceNum, //避免设置为0被误认为是不使用RBF的
		Sequence: sequence,
	}
}

func NewRBFActive() *RBFConfig {
	return NewRBFConfig(wire.MaxTxInSequenceNum - 2) // recommended sequence BTC推荐的默认启用RBF的就是这个数 // 选择 wire.MaxTxInSequenceNum - 2 而不是 wire.MaxTxInSequenceNum - 1 是出于一种谨慎性和规范性的考虑。虽然在技术上 wire.MaxTxInSequenceNum - 1 也可以支持 RBF，但 -2 更为常用
}

func NewRBFNotUse() *RBFConfig {
	return NewRBFConfig(wire.MaxTxInSequenceNum) //当两个元素都为零值时表示不启用RBF机制，当然这里设置为  wire.MaxTxInSequenceNum 也行，逻辑已经做过判定
}

func (cfg *RBFConfig) GetSequence() uint32 {
	if cfg.AllowRBF || cfg.Sequence > 0 { //启用RBF机制，精确的RBF逻辑
		return cfg.Sequence
	}
	return wire.MaxTxInSequenceNum //当两个元素都为零值时表示不启用RBF机制-因此这里使用默认的最大值表示不启用
}
