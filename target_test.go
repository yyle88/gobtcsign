package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

func TestNewAddressTuple_ValidAddress(t *testing.T) {
	// 测试有效地址
	address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"

	res := NewAddressTuple(address)
	t.Log("address:", res.Address)

	require.Equal(t, address, res.Address)
	require.NoError(t, res.VerifyMatch(&chaincfg.MainNetParams))

	pkScript, err := res.GetPkScript(&chaincfg.MainNetParams)
	require.NoError(t, err)
	t.Log("pk-script:", pkScript)

	expected := []byte{118, 169, 20, 98, 233, 7, 177, 92, 191, 39, 213, 66, 83, 153, 235, 246, 240, 251, 80, 235, 184, 143, 24, 136, 172}
	require.Equal(t, expected, pkScript)
}

func TestNewAddressTuple_AddressPredefined_PkScriptPredefined(t *testing.T) {
	// 测试预设 PkScript
	address := "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa"

	res := &AddressTuple{
		Address:  address,
		PkScript: []byte{118, 169, 20, 98, 233, 7, 177, 92, 191, 39, 213, 66, 83, 153, 235, 246, 240, 251, 80, 235, 184, 143, 24, 136, 172},
	}
	t.Log("address:", res.Address)
	t.Log("predefined pk-script:", res.PkScript)

	require.Equal(t, address, res.Address)
	require.NoError(t, res.VerifyMatch(&chaincfg.MainNetParams))

	pkScript, err := res.GetPkScript(&chaincfg.MainNetParams)
	require.NoError(t, err)

	expected := []byte{118, 169, 20, 98, 233, 7, 177, 92, 191, 39, 213, 66, 83, 153, 235, 246, 240, 251, 80, 235, 184, 143, 24, 136, 172}
	require.Equal(t, expected, pkScript)
}

func TestNewAddressTuple_PkScriptPredefined(t *testing.T) {
	// 测试预设 PkScript
	res := &AddressTuple{
		Address:  "",
		PkScript: []byte{118, 169, 20, 98, 233, 7, 177, 92, 191, 39, 213, 66, 83, 153, 235, 246, 240, 251, 80, 235, 184, 143, 24, 136, 172},
	}
	t.Log("predefined pk-script:", res.PkScript)

	require.NoError(t, res.VerifyMatch(&chaincfg.MainNetParams))

	pkScript, err := res.GetPkScript(&chaincfg.MainNetParams)
	require.NoError(t, err)

	expected := []byte{118, 169, 20, 98, 233, 7, 177, 92, 191, 39, 213, 66, 83, 153, 235, 246, 240, 251, 80, 235, 184, 143, 24, 136, 172}
	require.Equal(t, expected, pkScript)
}
