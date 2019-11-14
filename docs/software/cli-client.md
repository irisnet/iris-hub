# IRIS Command Line Client

## Introduction

`iriscli` is a client for the IRISnet network. IRISnet users can use `iriscli` to send different transactions and query the blockchain data.

## iriscli Work Directory

The default work directory for the `iriscli` is `$HOME/.iriscli`, which is mainly used to save configuration files and data. The IRISnet `key` data is saved in the work directory of `iriscli`. You can also specify the `iriscli`  work directory by `--home`.


## iriscli --node

The rpc address of the `iris` node. Transactions and query requests are sent to the process listening to this port. The default is `tcp://localhost:26657`, and the rpc address can also be specified by `--node`.

## iriscli config command

The `iriscli config` command interactively configures some default parameters, such as chain-id, home, fee, and node.

## Fee and Gas

`iriscli` can send a transaction with the fee specified by `--fee` and gas(the default is 50000) specified by `--gas` . The fee divided by the gas is gas price, and the gas price can't be less than the minimum value(6000 Nano) set by the blockchain. The remaining fees after the completion of the entire transaction will be returned to the user. You can set `--gas="simulate"`, which can estimate the gas consumed by the transaction through the simulation run, and multiply the coefficient specified by `--gas-adjustment`(default 1.5) to get the final gas as the transaction gas. Finally, the transaction will be broadcast.

```
iriscli bank send --amount=1iris --fee=0.3iris  --chain-id=<chain-id> --from=<user> --to=<address> --commit --gas="simulate"
```

## Dry-run Mode

`iriscli` turns off dry-run mode by default. If you want to open the dry-run mode, you can use `--dry-run`. It is similar to the simulation processing logic, and can calculate the gas that needs to be consumed, but then it will not broadcast to the full node, directly return and print the gas consumed this time.

Example: Send a command using dry-run mode

```
iriscli gov submit-proposal --title="test" --description="test" --type="Parameter" --deposit=600iris --param='mint/Inflation=0.050' --from=<user> --chain-id=<chain-id> --fee=0.3iris --dry-run
```

Print：

```
estimated gas = 18604
```

## Async Mode

async：Without any validation of the transaction, return the hash of the transaction immediately

sync：Verify the legality of the transaction (transaction format and signature), return the result and transaction hash. Transaction waiting to be packaged out in the blockchain

commit：Waiting for the transaction to be packaged in the blockchain before returning the complete execution result of the transaction，the request will be blocked until transaction return or timeout.

The default transaction mode of `iriscli` is `sync`. 

## Generate Only

`generate-only` is turned off by default, but `--generate-only` can be enabled to send the transaction, and then the unsigned transaction generated by the command line will be printed.

Example: Enable generate-only to generate unsigned transaction

```
iriscli gov submit-proposal --chain-id=<chain-id> --from=<user> --fee=0.3iris --description="test" --title="test" --usage="Burn" --percent=0.0000000001 --type="TxTaxUsage" --deposit=1000iris --generate-only
```

Print：

```json
{
  "msg": [
    {
      "type": "irishub/gov/MsgSubmitTxTaxUsageProposal",
      "value": {
        "MsgSubmitProposal": {
          "title": "test",
          "description": "test",
          "proposal_type": "CommunityTaxUsage",
          "proposer": "iaa1ljemm0yznz58qxxs8xyak7fashcfxf5lgl4zjx",
          "initial_deposit": [
            {
              "denom": "iris-atto",
              "amount": "1000000000000000000000"
            }
          ],
          "params": null
        },
        "usage": "Burn",
        "dest_address": "faa108w3x8",
        "percent": "0.0000000001"
      }
    }
  ],
  "fee": {
    "amount": [
      {
        "denom": "iris-atto",
        "amount": "300000000000000000"
      }
    ],
    "gas": "50000"
  },
  "signatures": null,
  "memo": ""
}
```

## Trust Node 

The trust-node is true by default. When the trust-node is true, the `iriscli` only queries the data and does not perform a Merkle-proof verification on the data. You can also use the `--trust-node=false` to perform Merkle-proof verification on the data obtained by the query.
