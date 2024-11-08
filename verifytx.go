package gobtcsign

import (
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcwallet/wallet/txauthor"
	"github.com/pkg/errors"
)

type VerifyTxInputParam struct {
	Sender AddressTuple
	Amount int64
}

func NewVerifyTxInputParam(senderAddress string, amount int64) *VerifyTxInputParam {
	return &VerifyTxInputParam{
		Sender: *NewAddressTuple(senderAddress),
		Amount: amount,
	}
}

func NewVerifyTxInputNotAmountParams(senders []string, netParams *chaincfg.Params) ([]*VerifyTxInputParam, error) {
	var results = make([]*VerifyTxInputParam, 0, len(senders))

	var a2pksMap = make(map[string][]byte, len(senders))
	for _, address := range senders {
		pkScript, ok := a2pksMap[address]
		if !ok {
			pks, err := GetAddressPkScript(address, netParams)
			if err != nil {
				return nil, errors.WithMessage(err, "cannot get pk-script")
			}
			a2pksMap[address] = pks
			pkScript = pks
		}

		results = append(results, &VerifyTxInputParam{
			Sender: AddressTuple{
				Address:  address,
				PkScript: pkScript,
			},
			Amount: 0, //绝大多数的签名，比如，P2PKH 签名，不将 amount 包含在生成的签名哈希中，因此也不验证它，随便填都行
		})
	}
	return results, nil
}

type VerifyTxInputsType struct {
	PkScripts [][]byte
	InAmounts []btcutil.Amount
}

func NewVerifyTxInputsType(inputList []*VerifyTxInputParam, netParams *chaincfg.Params) (*VerifyTxInputsType, error) {
	var pkScripts = make([][]byte, 0, len(inputList))
	var inAmounts = make([]btcutil.Amount, 0, len(inputList))
	for idx := range inputList {
		pkScript, err := inputList[idx].Sender.GetPkScript(netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong address->pk-script")
		}
		pkScripts = append(pkScripts, pkScript)
		inAmounts = append(inAmounts, btcutil.Amount(inputList[idx].Amount))
	}
	return &VerifyTxInputsType{
		PkScripts: pkScripts,
		InAmounts: inAmounts,
	}, nil
}

/*
VerifyP2PKHSign 验证签名是否有效，只有P2PKH的验证可以不验证数量，因此这里写个简易的函数，以便在需要的时候能够快速派上用场

这个函数的参数是：
  - 当前utxo持有者的地址，也就是这个utxo在谁的钱包里，这个数组元素可以重复，特别是在单签名的场景里。
  - 通常就是发送者的地址，假如使用两个utxo就是两个相同的发送者地址，假如是多签的情况，再看具体情况。

虽然 txid 和 vout 是 UTXO 的唯一标识，但在用户和系统交互的层面上，使用地址更为直观。
同时，验证过程中的安全性主要依赖于签名的有效性，而这种有效性是通过地址和相应的 pkScript 来实现的。
因此，系统选择通过utxo的来源地址来处理签名验证，而不是直接使用 txid 和 vout。

因此这里就是给的utxo的来源地址列表（按正确顺序排列，而且条数要相同）。
*/
func VerifyP2PKHSign(msgTx *wire.MsgTx, senders []string, netParams *chaincfg.Params) error {
	inputList, err := NewVerifyTxInputNotAmountParams(senders, netParams)
	if err != nil {
		return errors.WithMessage(err, "wrong new-input-params")
	}
	return VerifySignV2(msgTx, inputList, netParams)
}

/*
VerifySignV2 验证签名是否有效，同样的逻辑实现第二遍是为了简化参数，以便在需要的时候能够快速派上用场
验证签名的主要逻辑就是验证交易中的输入（vin）是否有效。
  - 签名的核心目的：确保交易中的 vin 确实有权花费引用的 UTXO。
  - 签名验证：主要是通过公钥和私钥的配对来验证引用的 UTXO 是否被合法使用。

在使用 SignP2PKH 时，如果你的交易是基于 P2PKH (Pay to Public Key Hash) 的传统交易，而不是 SegWit (BIP143) 交易，那么输入金额不会影响签名验证。
这是因为 P2PKH 签名不将 amount 包含在生成的签名哈希中。也就是说，amount 不会直接影响签名的生成和验证。
因此在这种情况下 amount 直接填 0 也行，填真实值也行。

在使用 SignP2WPKH 时，需要验证数量

逻辑中用到 NewSigCache 缓存功能
如果你的交易验证中有可能存在重复的 pkScript，那么使用 NewSigCache 来创建缓存是一个明智的选择，可以提高性能。
但如果你的场景非常简单且输入数量有限，设置为 nil 或 0 也完全可以接受。根据实际需求做出选择即可。

NewSigCache 创建的缓存通常不需要显式关闭或清理。它是一个内存中的数据结构，生命周期与其所在的应用程序或模块相同。
*/
func VerifySignV2(msgTx *wire.MsgTx, inputList []*VerifyTxInputParam, netParams *chaincfg.Params) error {
	inputsItem, err := NewVerifyTxInputsType(inputList, netParams)
	if err != nil {
		return errors.WithMessage(err, "wrong params-to-inputs")
	}
	return VerifySignV3(msgTx, inputsItem)
}

func VerifySignV3(msgTx *wire.MsgTx, inputsItem *VerifyTxInputsType) error {
	inputFetcher, err := txauthor.TXPrevOutFetcher(msgTx, inputsItem.PkScripts, inputsItem.InAmounts)
	if err != nil {
		return errors.WithMessage(err, "wrong cannot-create-pre-out-cache")
	}
	sigHashCache := txscript.NewTxSigHashes(msgTx, inputFetcher)

	inputOuts := NewInputOutsV2(inputsItem.PkScripts, inputsItem.InAmounts)

	return VerifySign(msgTx, inputOuts, inputFetcher, sigHashCache)
}
