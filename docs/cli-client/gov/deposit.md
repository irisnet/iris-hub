# iriscli gov deposit

## Description
 
Deposit tokens for activing proposal
 
## Usage
 
```
iriscli gov deposit [flags]
```

Print help messages:

```
iriscli gov deposit --help
```
## Flags
 
| Name, shorthand  | Default                    | Description                                                                                                                                          | Required |
| ---------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- | -------- |
| --deposit        |                            | [string] Deposit of proposal                                                                                                                         | Yes      |
| --proposal-id    |                            | [string] ProposalID of proposal depositing on                                                                                                        | Yes      |
## Examples

### Deposit

```shell
iriscli gov deposit --chain-id=test --proposal-id=1 --deposit=50iris --from=node0 --fee=0.01iris
```

After you enter the correct password, you could deposit 50iris to make your proposal active which can be voted, after you enter the correct password, you're done with depositing iris tokens for an activing proposal.

```txt
Password to sign with 'node0':
Committed at block 473 (tx hash: 0309E969589F19A9D9E4BFC9479720487FBF929ED6A88824414C5E7E91709206, response: {Code:0 Data:[] Log:Msg 0:  Info: GasWanted:200000 GasUsed:6710 Tags:[{Key:[97 99 116 105 111 110] Value:[100 101 112 111 115 105 116] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[100 101 112 111 115 105 116 101 114] Value:[102 97 97 49 52 113 53 114 102 57 115 108 50 100 113 100 50 117 120 114 120 121 107 97 102 120 113 51 110 117 51 108 106 50 102 112 57 108 55 112 103 100] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[112 114 111 112 111 115 97 108 45 105 100] Value:[49] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0} {Key:[99 111 109 112 108 101 116 101 67 111 110 115 117 109 101 100 84 120 70 101 101 45 105 114 105 115 45 97 116 116 111] Value:[34 51 51 53 53 48 48 48 48 48 48 48 48 48 48 48 34] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}] Codespace: XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0})
{
   "tags": {
     "action": "deposit",
     "completeConsumedTxFee-iris-atto": "\"335500000000000\"",
     "depositor": "faa14q5rf9sl2dqd2uxrxykafxq3nu3lj2fp9l7pgd",
     "proposal-id": "1"
   }
 }
```
