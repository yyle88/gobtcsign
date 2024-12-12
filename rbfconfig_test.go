package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
)

func TestRBFConfig_NewRBFConfig_AllowRBFTrue(t *testing.T) {
	cfg := NewRBFConfig(wire.MaxTxInSequenceNum - 1)
	require.Equal(t, true, cfg.AllowRBF)
	require.Equal(t, wire.MaxTxInSequenceNum-1, cfg.Sequence)
}

func TestRBFConfig_NewRBFConfig_AllowRBFFalse(t *testing.T) {
	cfg := NewRBFConfig(wire.MaxTxInSequenceNum)
	require.Equal(t, false, cfg.AllowRBF)
	require.Equal(t, wire.MaxTxInSequenceNum, cfg.Sequence)
}

func TestRBFConfig_NewRBFActive(t *testing.T) {
	cfg := NewRBFActive()
	require.Equal(t, true, cfg.AllowRBF)
	require.Equal(t, wire.MaxTxInSequenceNum-2, cfg.Sequence)
}

func TestRBFConfig_NewRBFNotUse(t *testing.T) {
	cfg := NewRBFNotUse()
	require.Equal(t, false, cfg.AllowRBF)
	require.Equal(t, wire.MaxTxInSequenceNum, cfg.Sequence)
}

func TestRBFConfig_GetSequence_AllowRBFTrue(t *testing.T) {
	cfg := &RBFConfig{AllowRBF: true, Sequence: wire.MaxTxInSequenceNum - 1}
	require.Equal(t, wire.MaxTxInSequenceNum-1, cfg.GetSequence())
}

func TestRBFConfig_GetSequence_AllowRBFFalseWithCustomSequence(t *testing.T) {
	cfg := &RBFConfig{AllowRBF: false, Sequence: wire.MaxTxInSequenceNum - 2}
	require.Equal(t, wire.MaxTxInSequenceNum-2, cfg.GetSequence())
}

func TestRBFConfig_GetSequence_AllowRBFFalseWithDefaultSequence(t *testing.T) {
	cfg := &RBFConfig{AllowRBF: false, Sequence: 0}
	require.Equal(t, wire.MaxTxInSequenceNum, cfg.GetSequence())
}

func TestRBFConfig_NewRBFConfig_AllowRBFTrueWithZeroSequence(t *testing.T) {
	// 测试当 AllowRBF 为 true 且 Sequence 为 0 时的行为
	cfg := NewRBFConfig(0)
	require.Equal(t, true, cfg.AllowRBF)      // AllowRBF 应该是 true
	require.Equal(t, uint32(0), cfg.Sequence) // Sequence 应该是 0
}

func TestRBFConfig_NewRBFConfig_WithLargeSequence(t *testing.T) {
	// 测试 Sequence 设置为一个较大的值
	cfg := NewRBFConfig(999999999)
	require.Equal(t, true, cfg.AllowRBF)              // AllowRBF 应该是 true
	require.Equal(t, uint32(999999999), cfg.Sequence) // Sequence 应该是设置的值
}

func TestRBFConfig_GetSequence_AllowRBFTrueWithZeroSequence(t *testing.T) {
	// 测试当 AllowRBF 为 true 且 Sequence 为 0 时，GetSequence 的返回值
	cfg := &RBFConfig{AllowRBF: true, Sequence: 0}
	require.Equal(t, uint32(0), cfg.GetSequence()) // 如果 AllowRBF 为 true，应该返回 Sequence 本身的值，即 0
}

func TestRBFConfig_GetSequence_AllowRBFTrueWithLargeSequence(t *testing.T) {
	// 测试当 AllowRBF 为 true 且 Sequence 设置为较大值时，GetSequence 的返回值
	cfg := &RBFConfig{AllowRBF: true, Sequence: 999999999}
	require.Equal(t, uint32(999999999), cfg.GetSequence()) // 应该返回设置的值
}

func TestRBFConfig_GetSequence_AllowRBFFalseWithSequenceEqualToMaxTxInSequenceNum(t *testing.T) {
	// 测试当 Sequence 设置为 wire.MaxTxInSequenceNum 且 AllowRBF 为 false 时
	cfg := &RBFConfig{AllowRBF: false, Sequence: wire.MaxTxInSequenceNum}
	require.Equal(t, wire.MaxTxInSequenceNum, cfg.GetSequence()) // 由于 Sequence 为 MaxTxInSequenceNum，应该返回 MaxTxInSequenceNum
}
