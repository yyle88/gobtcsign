package gosignbtc

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
)

func TestCreateWalletP2PKH(t *testing.T) {
	netParams := chaincfg.MainNetParams

	address, private, err := CreateWalletP2PKH(&netParams)
	require.NoError(t, err)
	t.Log(address)
	t.Log(private)
}

func TestCreateWalletP2PKH_2(t *testing.T) {
	netParams := chaincfg.TestNet3Params

	address, private, err := CreateWalletP2PKH(&netParams)
	require.NoError(t, err)
	t.Log(address)
	t.Log(private)
}
