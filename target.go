package gobtcsign

import (
	"bytes"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pkg/errors"
)

type AddressTuple struct {
	Address  string //钱包地址 和 公钥脚本 二选一填写即可
	PkScript []byte //公钥脚本 和 钱包地址 二选一填写即可 PkScript（Public Key Script）在拼装交易和签名时使用
}

func NewAddressTuple(address string) *AddressTuple {
	return &AddressTuple{
		Address:  address,
		PkScript: nil, //这里 address 和 pk-script 是二选一的，因此不设，在后续的逻辑里会根据地址获得 pk-script 信息
	}
}

// GetPkScript 获得公钥文本，当公钥文本存在时就用已有的，否则就根据地址计算
func (one *AddressTuple) GetPkScript(netParams *chaincfg.Params) ([]byte, error) {
	if len(one.PkScript) > 0 && len(one.Address) > 0 {
		// 这里的目的不是缓存而是两个参数都可以填，但当两个参数都填的时候就得保证匹配，避免出问题
		pkScript, err := GetAddressPkScript(one.Address, netParams)
		if err != nil {
			return nil, errors.WithMessage(err, "wrong-address")
		}
		if bytes.Compare(one.PkScript, pkScript) != 0 {
			return nil, errors.New("address-pk-script-mismatch")
		}
		return pkScript, nil
	}
	if len(one.PkScript) > 0 {
		return one.PkScript, nil //假如有就直接返回，否则就根据地址计算
	}
	if one.Address != "" {
		return GetAddressPkScript(one.Address, netParams) //这里不用做缓存避免增加复杂度
	}
	return nil, errors.New("no-pk-script-no-address")
}

func (one *AddressTuple) VerifyMatch(netParams *chaincfg.Params) error {
	if one.Address != "" && len(one.PkScript) > 0 {
		pkScript, err := GetAddressPkScript(one.Address, netParams)
		if err != nil {
			return errors.WithMessage(err, "wrong-address")
		}
		if bytes.Compare(one.PkScript, pkScript) != 0 {
			return errors.New("address-pk-script-mismatch")
		}
	}
	return nil
}
