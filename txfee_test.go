package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/wire"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestCalculateMsgTxSize(t *testing.T) {
	const senderAddress = "nVnVaL5e4L2GDRha9aQ7KiSXDnqjUUz1K4"

	netParams := dogecoin.TestNetParams

	param := &CustomParam{
		VinList: []VinType{
			{
				OutPoint: *MustNewOutPoint("5ae74f2d6c4a0513e3c75484a726820c2b0653c2b26352afe97f4bf813dcf859", 0),
				Sender:   *NewAddressTuple(senderAddress),
				Amount:   1000000,
				RBFInfo:  *NewRBFNotUse(),
			},
			{
				OutPoint: *MustNewOutPoint("336d48ad5b7f2c72b98adc19cd7a56083f8e52f87958368810d47354b97acb38", 0),
				Sender:   *NewAddressTuple(senderAddress),
				Amount:   1000000,
				RBFInfo:  *NewRBFNotUse(),
			},
			{
				OutPoint: *MustNewOutPoint("e4ab15e75aa66fb67b02a10bf2772269d1b6b135ebbd6480f29eef2fc8825934", 0),
				Sender:   *NewAddressTuple(senderAddress),
				Amount:   1000000,
				RBFInfo:  *NewRBFNotUse(),
			},
		},
		OutList: []OutType{
			{
				Target: *NewAddressTuple("nhrZGEEh7JgVV3T1ncnUdTDZsByNnkmipc"),
				Amount: 1000000,
			},
		},
		RBFInfo: *NewRBFActive(),
	}

	signParam, err := param.GetSignParam(&netParams)
	require.NoError(t, err)

	changeAddress, err := btcutil.DecodeAddress(senderAddress, &netParams)
	require.NoError(t, err)

	size, err := CalculateMsgTxSize(signParam.MsgTx, changeAddress)
	require.NoError(t, err)

	// see doge testnet tx 8f55eb7057c6fa524dff88b0bfa0e208bec3db159c374ae3ff29889c9e4d33dd
	t.Log(size) // almost same size with the chain size
	require.Equal(t, 525, size)

	txFee, err := CalculateMsgTxFee(signParam.MsgTx, changeAddress, 1500000, dogecoin.NewDogeDustFee())
	require.NoError(t, err)

	t.Log(txFee) //这里打印出来单位是BTC，但实际是DOGE，但是不用在意这些细节
	require.Equal(t, btcutil.Amount(787500), txFee)

	changeAmount := param.GetFee() - txFee
	t.Log(changeAmount)
	require.Equal(t, btcutil.Amount(1212500), changeAmount)

	// append change out to outputs
	signParam.MsgTx.AddTxOut(wire.NewTxOut(int64(changeAmount), MustGetPkScript(changeAddress)))
}
