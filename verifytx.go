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

/*
VerifyP2PKHSignV2 验证签名是否有效，同样的逻辑实现第二遍是为了简化参数，以便在需要的时候能够快速派上用场
验证签名的主要逻辑就是验证交易中的输入（vin）是否有效。
  - 签名的核心目的：确保交易中的 vin 确实有权花费引用的 UTXO。
  - 签名验证：主要是通过公钥和私钥的配对来验证引用的 UTXO 是否被合法使用。

在使用 SignP2PKH 时，如果你的交易是基于 P2PKH (Pay to Public Key Hash) 的传统交易，而不是 SegWit (BIP143) 交易，那么输入金额不会影响签名验证。
这是因为 P2PKH 签名不将 amount 包含在生成的签名哈希中。也就是说，amount 不会直接影响签名的生成和验证。
因此在这种情况下 amount 直接填 0 也行，填真实值也行。

逻辑中用到 NewSigCache 缓存功能
如果你的交易验证中有可能存在重复的 pkScript，那么使用 NewSigCache 来创建缓存是一个明智的选择，可以提高性能。
但如果你的场景非常简单且输入数量有限，设置为 nil 或 0 也完全可以接受。根据实际需求做出选择即可。

NewSigCache 创建的缓存通常不需要显式关闭或清理。它是一个内存中的数据结构，生命周期与其所在的应用程序或模块相同。
*/
func VerifyP2PKHSignV2(msgTx *wire.MsgTx, inputList []*VerifyTxInputParam, netParams *chaincfg.Params) error {
	var prevScripts = make([][]byte, 0, len(inputList))
	var inputValues = make([]btcutil.Amount, 0, len(inputList))
	for idx := range inputList {
		pkScript, err := inputList[idx].Sender.GetPkScript(netParams)
		if err != nil {
			return errors.WithMessage(err, "wrong address->pk-script")
		}
		prevScripts = append(prevScripts, pkScript)
		inputValues = append(inputValues, btcutil.Amount(inputList[idx].Amount))
	}

	inputFetcher, err := txauthor.TXPrevOutFetcher(msgTx, prevScripts, inputValues)
	if err != nil {
		return errors.WithMessage(err, "cannot create prev out cache")
	}
	sigHashCache := txscript.NewTxSigHashes(msgTx, inputFetcher)

	sigCache := txscript.NewSigCache(uint(len(inputList))) //设置为输入的长度是较好的，当然，更大量的计算时也可使用全局的cache

	for txIdx, prevScript := range prevScripts {
		vm, err := txscript.NewEngine(
			prevScript,
			msgTx,
			txIdx,
			txscript.StandardVerifyFlags,
			sigCache,
			sigHashCache,
			int64(inputValues[txIdx]),
			inputFetcher,
		)
		if err != nil {
			return errors.WithMessage(err, "cannot create script engine")
		}
		if err = vm.Execute(); err != nil {
			return errors.WithMessage(err, "cannot validate transaction")
		}
	}
	return nil
}

/*
VerifyP2PKHSignV3 验证签名是否有效，同样的逻辑实现第三遍是为了简化参数，以便在需要的时候能够快速派上用场

这个函数的参数是：
  - 当前utxo持有者的地址，也就是这个utxo在谁的钱包里，这个数组元素可以重复，特别是在单签名的场景里。
  - 通常就是发送者的地址，假如使用两个utxo就是两个相同的发送者地址，假如是多签的情况，再看具体情况。

虽然 txid 和 vout 是 UTXO 的唯一标识，但在用户和系统交互的层面上，使用地址更为直观。
同时，验证过程中的安全性主要依赖于签名的有效性，而这种有效性是通过地址和相应的 pkScript 来实现的。
因此，系统选择通过utxo的来源地址来处理签名验证，而不是直接使用 txid 和 vout。

因此这里就是给的utxo的来源地址列表（按正确顺序排列，而且条数要相同）。
*/
func VerifyP2PKHSignV3(msgTx *wire.MsgTx, sendersAddresses []string, netParams *chaincfg.Params) error {
	var inputList = make([]*VerifyTxInputParam, 0, len(sendersAddresses))

	var address2pkMap = make(map[string][]byte, len(sendersAddresses))
	for _, address := range sendersAddresses {
		pkScriptValue, ok := address2pkMap[address]
		if !ok {
			pkScript, err := GetAddressPkScript(address, netParams)
			if err != nil {
				return errors.WithMessage(err, "cannot get pk-script")
			}
			address2pkMap[address] = pkScript
			pkScriptValue = pkScript
		}

		inputList = append(inputList, &VerifyTxInputParam{
			Sender: AddressTuple{
				Address:  address,
				PkScript: pkScriptValue,
			},
			Amount: 0, //绝大多数的签名，比如，P2PKH 签名，不将 amount 包含在生成的签名哈希中，因此也不验证它，随便填都行
		})
	}
	return VerifyP2PKHSignV2(msgTx, inputList, netParams)
}
