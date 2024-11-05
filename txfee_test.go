package gobtcsign

import (
	"testing"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/require"
	"github.com/yyle88/gobtcsign/dogecoin"
)

func TestEstimateTxSize(t *testing.T) {
	{
		const txHex = "020000000001017fb11165fc5a6edf3bf06176c8915b22e20c3c48966d1d5c5673d4bc76a98c6e0100000000fdffffff024a140000000000001976a9145e3a929c6f941ad02c352b47d33a65bd160afe2f88accc8c0700000000001600149aba0ec86126b94b4dfdff48c4855ec49975c8e50247304402200615b90428b7e857d074ef17da69d5d533adc2c02fb739b04ecb88e38723afcd0220070f30781a8020ba0330a96b2a58358c5b5073ed391cee50ea1650e3972fcb0d012102030df0337a88e6f5c77f593eeb8b9d425742fbee78a4dcfeb66f90dec8e30bf25d560b00"

		mstTx, err := NewMsgTxFromHex(txHex)
		require.NoError(t, err)

		txHash := GetTxHash(mstTx)
		t.Log(txHash)
		require.Equal(t, "fa3467452918627ebd63a3e8570e70d38b0eefef683347510b204ba6962ebe44", txHash)

		size := GetMsgTxVSize(mstTx)
		t.Log(size)
		require.Equal(t, 144, size)
	}

	netParams := chaincfg.MainNetParams

	param := &CustomParam{
		VinList: []VinType{
			{
				OutPoint: *MustNewOutPoint("6e8ca976bcd473565c1d6d96483c0ce2225b91c87661f03bdf6e5afc6511b17f", 1),
				Sender:   *NewAddressTuple("bc1q7klfa2srcuwfeen625t2ydxpfswuqns4787w3v"),
				Amount:   500134,
				RBFInfo:  *NewRBFNotUse(),
			},
		},
		OutList: []OutType{
			{
				Target: *NewAddressTuple("19bEgjc6Cuu27VVF8cmN4bE5ykETqyZeXF"),
				Amount: 5194,
			},
			{
				Target: *NewAddressTuple("bc1qn2aqajrpy6u5kn0alayvfp27cjvhtj89mqfkhd"),
				Amount: 494796,
			},
		},
		RBFInfo: *NewRBFActive(),
	}

	size, err := EstimateTxSize(param, &netParams, NewNoChange())
	require.NoError(t, err)

	t.Log("estimate-tx-size:", size)
	require.Equal(t, 144, size)

	t.Log(param.GetFee())
}

func TestEstimateTxSize_VIN_1_P2PKH(t *testing.T) {
	{
		const txHex = "0100000001a4c91c9720157a5ee582a7966471d9c70d0a860fa7757b4c42a535a12054a4c9000000006c493046022100d49c452a00e5b1213ac84d92269510a05a584a4d0949bd7d0ad4e3408ac8e80a022100bf98707ffaf1eb9dff146f7da54e68651c0a27e3653ec3882b7a95202328579c01210332d98672a4246fe917b9c724c339e757d46b1ffde3fb27fdc680b4bb29b6ad59ffffffff02a0860100000000001976a9144fb55ee0524076acd4c14e7773561e4c298c8e2788ac20688a0b000000001976a914cb7f6bb8e95a2cd06423932cfbbce73d16a18df088ac00000000"

		mstTx, err := NewMsgTxFromHex(txHex)
		require.NoError(t, err)

		txHash := GetTxHash(mstTx)
		t.Log(txHash)
		require.Equal(t, "1d07ae04c8114064b941b3e7acd550a55ebfac81a1008d63f4a456b59bdda680", txHash)

		size := GetMsgTxVSize(mstTx)
		t.Log(size)
		require.Equal(t, 227, size)
	}

	netParams := chaincfg.MainNetParams

	param := &CustomParam{
		VinList: []VinType{
			{
				OutPoint: *MustNewOutPoint("c9a45420a135a5424c7b75a70f860a0dc7d9716496a782e55e7a1520971cc9a4", 0),
				Sender:   *NewAddressTuple("16bXaHE1R8Vgm8RBqMHhfcDcZ3msPd5K9R"),
				Amount:   193730000,
				RBFInfo:  *NewRBFNotUse(),
			},
		},
		OutList: []OutType{
			{
				Target: *NewAddressTuple("18GTf1KdLt9ihs9tfo3v142uxWrku8E71t"),
				Amount: 100000,
			},
			{
				Target: *NewAddressTuple("1KYzpJ1e2soWMNQErczu9sweyvnKtUXbXN"),
				Amount: 193620000,
			},
		},
		RBFInfo: *NewRBFActive(),
	}

	size, err := EstimateTxSize(param, &netParams, NewNoChange())
	require.NoError(t, err)

	t.Log("estimate-tx-size:", size)
	require.Equal(t, 227, size)

	t.Log(param.GetFee())
}

func TestEstimateTxFee(t *testing.T) {
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

	changeAddress, err := btcutil.DecodeAddress("nVnVaL5e4L2GDRha9aQ7KiSXDnqjUUz1K4", &netParams)
	require.NoError(t, err)

	size, err := EstimateTxSize(param, &netParams, &ChangeTo{AddressX: changeAddress})
	require.NoError(t, err)

	// see doge testnet tx 8f55eb7057c6fa524dff88b0bfa0e208bec3db159c374ae3ff29889c9e4d33dd
	t.Log("estimate-tx-size:", size) // almost same size with the chain size //这是预估值 略微 >= 实际值
	require.Equal(t, 525, size)

	txFee, err := EstimateTxFee(param, &netParams, &ChangeTo{AddressX: changeAddress}, 1500000, dogecoin.NewDogeDustFee())
	require.NoError(t, err)

	t.Log(txFee) //这里打印出来单位是BTC，但实际是DOGE，但是不用在意这些细节
	require.Equal(t, btcutil.Amount(787500), txFee)

	changeAmount := param.GetChangeAmountWithFee(txFee)
	t.Log(changeAmount)
	require.Equal(t, btcutil.Amount(1212500), changeAmount)

	// append change out to outputs
	param.OutList = append(param.OutList, OutType{
		Target: *NewAddressTuple("nVnVaL5e4L2GDRha9aQ7KiSXDnqjUUz1K4"),
		Amount: int64(changeAmount),
	})
	// make sure the tx fee is same
	require.Equal(t, btcutil.Amount(787500), param.GetFee())
}
