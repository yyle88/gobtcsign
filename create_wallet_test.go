package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestCreateWalletP2PKH_BTC(t *testing.T) {
	netParams := chaincfg.MainNetParams

	address, private, err := CreateWalletP2PKH(&netParams)
	require.NoError(t, err)
	t.Log(address)
	t.Log(private)
	t.Log(netParams.Name)
}

func TestCreateWalletP2PKH_BTC_testnet(t *testing.T) {
	netParams := chaincfg.TestNet3Params

	address, private, err := CreateWalletP2PKH(&netParams)
	require.NoError(t, err)
	t.Log(address)
	t.Log(private)
	t.Log(netParams.Name)
}

func TestCreateWalletP2PKH_DOGE(t *testing.T) {
	netParams := dogecoin.MainNetParams

	address, private, err := CreateWalletP2PKH(&netParams)
	require.NoError(t, err)
	t.Log(address)
	t.Log(private)
	t.Log(netParams.Name)
}

func TestCreateWalletP2PKH_DOGE_testnet(t *testing.T) {
	netParams := dogecoin.TestNetParams

	address, private, err := CreateWalletP2PKH(&netParams)
	require.NoError(t, err)
	t.Log(address)
	t.Log(private)
	t.Log(netParams.Name)
}

func TestCreateWalletP2WPKH_BTC(t *testing.T) {
	netParams := chaincfg.MainNetParams

	address, private, err := CreateWalletP2WPKH(&netParams)
	require.NoError(t, err)
	t.Log(address)
	t.Log(private)
	t.Log(netParams.Name)
}

func TestCreateWalletP2WPKH_BTC_testnet(t *testing.T) {
	netParams := chaincfg.TestNet3Params

	address, private, err := CreateWalletP2WPKH(&netParams)
	require.NoError(t, err)
	t.Log(address)
	t.Log(private)
	t.Log(netParams.Name)
}
