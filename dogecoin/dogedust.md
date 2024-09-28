# dogecoin dust fee estimate example

```go
maxSignedSize, _ := utils.BtcEstimatedTxSize(
    p2pkh, p2tr, p2wpkh, nested, outputs, changeAddress,
)
// 计算全局手续费
maxRequiredFee := txrules.FeeForSerializeSize(feeRatePerKb, maxSignedSize) + utils.GetSoftDustsFee(outputs)
// 手续费大于了除输出后的金额
// 尝试添加更多的utxo
if remainingAmount := inputAmount - targetAmount; remainingAmount < maxRequiredFee {
    targetFee = maxRequiredFee
    continue
}
// 假如有灰尘找零就提高预估的手续费
// 这是暂不舍掉灰尘找零的vout，而是补更多的INPUT把灰尘找零变为非灰尘找零，最后实质上并不会消耗更高的手续费
// 详见交易 TESTNET 的 HASH: 8f55eb7057c6fa524dff88b0bfa0e208bec3db159c374ae3ff29889c9e4d33dd 你就会秒懂啦
if changeAmount := inputAmount - targetAmount - maxRequiredFee; changeAmount >= constant.ChainMinDustOutput {
    //当这个数比刚性灰尘限额更大些（即能够提交给节点，暂不舍弃它），但又小于柔性限额时（即需要更高的手续费）
    if changeAmount < constant.ChainSoftDustLimit {
        //根据官方文档就需要再补些手续费
        maxRequiredFee += constant.ChainExtraDustsFee
        //这里是必然为真的，通过算术能够推断出来，但依然写在这里以免出现错误
        if remainingAmount := inputAmount - targetAmount; remainingAmount < maxRequiredFee {
            targetFee = maxRequiredFee
            continue
        }
    }
}
```

详情请看这个测试链的交易: 8f55eb7057c6fa524dff88b0bfa0e208bec3db159c374ae3ff29889c9e4d33dd
通过增补1个vin让找零变大而不再是灰尘就能解决问题

当然归根结底这个逻辑能生效是因为 constant.ChainSoftDustLimit == constant.ChainExtraDustsFee
假如不是这样的话可能还真得多给些手续费，因此这里的技巧需要掌握住
