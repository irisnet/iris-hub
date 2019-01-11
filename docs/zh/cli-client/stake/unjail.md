# iriscli stake unjail

## 介绍


在PoS网络中，验证人的收益主要来自于staking抵押获利，但是若他们不能保持在线，就会被当作一种作恶行为。系统会剥夺它作为验证人参与共识的资格。这样一来，它的状态会变成jailed，他们的投票权将立刻变为零。这种状态降持续一段时间。当jailed期结束，验证人节点的operator需要执行 unjail操作来让节点的状态变为unjailed，再次成为候选验证人。


## 用法

```
iriscli stake unjail [flags]
```

打印帮助信息

```
iriscli stake unjail --help
```

## 例子

### Unjail验证人节点

```
iriscli stake unjail --from=<key name> --fee=0.004iris --chain-id=<chain-id>
```
### 常见问题

* 检查这个验证人保持`jail`状态的截止时间：

```$xslt
iriscli stake signing-info fvp1zcjduepqewwc93xwvt0ym6prxx9ppfzeufs33flkcpu23n5eutjgnnqmgazsw54sfv --node=localhost:36657 --trust-node
```

如果此验证人状态为`jailed`，那么你可以看到它的jail状态的截止时间

```
Start height: 565, index offset: 2, jailed until: 2018-12-12 06:46:37.274910287 +0000 UTC, missed blocks counter: 2
```

如果你在jail状态的截止时间前执行`unjail` 命令，你会看到以下错误：

```$xslt
ERROR: Msg 0 failed: {"codespace":10,"code":102,"abci_code":655462,"message":"validator still jailed, cannot yet be unjailed"}
```

过了jail状态的截止时间后，你可以发送一个 `unjail` 交易. 

```
iriscli stake unjail --from=<key name> --fee=0.004iris --chain-id=test-irishub
```

输出:
```txt
Committed at block 306 (tx hash: 5A4C6E00F4F6BF795EB05D2D388CBA0E8A6E6CF17669314B1EE6A31729A22450, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:3398 Tags:[{Key:[97 99 116 105 111 110] Value:[115 101 114 118 105 99 101 45 119 105 116 104 100 114 97 119 45 102 101 101 115] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 54 55 57 54 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
```

```json
{
   "tags": {
     "action": "unjail",
     "completeConsumedTxFee-iris-atto": "\"918600000000000\"",
     "validator": "fva12zgt9hc5r5mnxegam9evjspgwhkgn4wzjxkvqy"
   }
 }
```
