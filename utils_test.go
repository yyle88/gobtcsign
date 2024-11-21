package gobtcsign

import (
	"encoding/base64"
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestGetAddressPkScript(t *testing.T) {
	netParams := dogecoin.TestNetParams
	pkScript := caseGetAddressPkScript(t, "nXMSrjEQXUJ77TQSeErpJMySy3kfSfwSCP", &netParams)

	// 这里写个简单的比较逻辑
	pkTarget, err := base64.StdEncoding.DecodeString("dqkUIqn5GrQ9r5dmyzOQ/RTb+qkqVQqIrA==")
	require.NoError(t, err)
	require.Equal(t, pkTarget, pkScript)
}

func caseGetAddressPkScript(t *testing.T, rawAddress string, netParams *chaincfg.Params) []byte {
	pkScript, err := GetAddressPkScript(rawAddress, netParams)
	require.NoError(t, err)
	return pkScript
}

func TestNewInputOuts(t *testing.T) {
	netParams := chaincfg.MainNetParams
	pkScript, err := GetAddressPkScript("tb1qvg2jksxckt96cdv9g8v9psreaggdzsrlm6arap", &netParams)
	require.NoError(t, err)
	outs := NewInputOuts([][]byte{pkScript, pkScript}, []int64{1234, 5678})
	require.Equal(t, 2, len(outs))
}

func TestGetMsgTxVSize(t *testing.T) {
	const txHex = "02000000000101fc74ccf02694e7fc82c6ee010a7d48b186a3b8503aa2f523112c12b83566fd770100000000fdffffff02241300000000000016001462152b40d8b2cbac358541d850c079ea10d1407f8596df020000000016001419372a643baa41e746fae4111b57e5c59f2d3315014011f36398ae04e0690ab0c47a4aed06425bb5b41ebb2349af8eccf4c1bf4c18ac498922549d977f83501bcd916d466f47f65cd1bde22cb018c5ffbc2770d69b9e6ac53000"

	mstTx, err := NewMsgTxFromHex(txHex)
	require.NoError(t, err)

	txHash := GetTxHash(mstTx)
	t.Log(txHash)
	require.Equal(t, "fb87cc4010bd4a34cb4be86f37182fada63c9923ae8eae5d2f793cb5f50c6328", txHash)

	size := GetMsgTxVSize(mstTx)
	t.Log(size)
	require.Equal(t, 130, size)
}
