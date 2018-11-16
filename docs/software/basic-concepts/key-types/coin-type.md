# Coin_Type

##  Definitions

Coin_type defines the available units of a kind of token in IRISnet. The developers can specify different coin_type for  their tokens. The native token in IRIShub is `iris`, which has following available units: `iris-milli`, `iris-micro`, `iris-nano`, `iris-pico`, `iris-femto` and `iris-atto`. The conversion relationship between them are as follows:

```
1 iris = 10^3 iris-milli
1 iris = 10^6 iris-micro
1 iris = 10^9 iris-nano
1 iris = 10^12 iris-pico
1 iris = 10^15 iris-femto
1 iris = 10^18 iris-atto
```

All the registered types of `iris` in the system can be used with transactions.

## Data Structure of coin_type

```golang
type CoinType struct {
	Name    string `json:"name"`
	MinUnit Unit   `json:"min_unit"`
	Units   Units  `json:"units"`
	Origin  Origin `json:"origin"`
	Desc    string `json:"desc"`
}
```

## Structure definition of Unit

```golang
type Unit struct {
	Denom   string `json:"denom"`
	Decimal int    `json:"decimal"`
}
```

* Name : The name of a token, which is also its default unit；for instance,the default unit of `iris` is `iris`.
* MinUnit：The  minimum unit of coin_type. 

The tokens in the system are all stored in the form of minimum unit, 
such as `iris-atto`. You could choose to use the minimum unit of the tokens when sending a transaction to the IRIShub. 
If you use the command line client, aka `iriscli`, you can use any system-recognized unit and the system 
will automatically convert to the minimum unit of this corresponding token. For example, if you execute `send`command 
to transfer 1iris, the command line will be processed as 10^18 iris-attos in the backend, and you will only 
see 10^18 `iris-attos` when searching the transaction details by transaction hash.



`Denom` is defined as the name of this unit, and `Decimal` is defined as the precision of the unit. 

For example, the precision of iris-atto is 18.

* `Unit` defines a set of units available under coin_type.
* `Origin` defines the source of the coin_type, with the value `Native` (inner system, iris for IRIShub), 
`External` (external system, such as eth for Ethereum, etc.), and `UserIssued` (user-defined).
* `Desc`：Description of the coin_type.

## Query of coin_type

If you want to query the coin_type configuration of a certain token, you can use the following command:

```golang
iriscli bank coin-type [coin_name]
```

If you query the `coin-type` of `iris` with `iriscli bank coin-type iris`
 
Example output:
```$xslt
{
  "name": "iris",
  "min_unit": {
    "denom": "iris-atto",
    "decimal": "18"
  },
  "units": [
    {
      "denom": "iris",
      "decimal": "0"
    },
    {
      "denom": "iris-milli",
      "decimal": "3"
    },
    {
      "denom": "iris-micro",
      "decimal": "6"
    },
    {
      "denom": "iris-nano",
      "decimal": "9"
    },
    {
      "denom": "iris-pico",
      "decimal": "12"
    },
    {
      "denom": "iris-femto",
      "decimal": "15"
    },
    {
      "denom": "iris-atto",
      "decimal": "18"
    }
  ],
  "origin": 1,
  "desc": "IRIS Network"
}
```