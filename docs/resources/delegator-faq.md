# Delegators

## What is a delegator?
People that cannot, or do not want to run validator operations, can still participate in the staking process as delegators. Indeed, validators are not chosen based on their own stake 
but based on their total stake, which is the sum of their own stake and of the stake that is delegated to them. This is an important property, as it makes delegators a safeguard against
 validators that exhibit bad behavior. If a validator misbehaves, its delegators will move their IRIS tokens away from it, thereby reducing its stake. Eventually, if a validator's stake falls 
 under the top 100 addresses with highest stake, it will exit the validator set.

## States for a Delegator

Delegators have the same state as their validator.

Note that delegation are not necessarily bonded. Tokens of each delegator can be delegated and bonded, delegated and unbonding, delegated and unbonded, or loose. 

## Common operation for Delegators

* Delegation

To delegate some IRIS token to a validator, you could run the following command:
```$xslt
iriscli stake delegate --address-delegator=<address-delegator> --address-validator=<address-validator> --chain-id=irishub --from=name --gas=50000 --fee=0.3iris --amount=10iris
```
> Please notice that the amount is under unit iris-atto, 1iris=10^18 iris-atto

* Query Delegation

You could query your delegation amount with the following command:

```$xslt
iriscli stake delegation --address-delegator=<address-delegator> --address-validator=<address-validator> --chain-id=irishub
```

The example output is the following:
```$xslt
Delegation
Delegator: iaa1je9qyff4qate4e0kthum0p8v7q7z8lr7phygv8
Validator: iaa1dmp6eyjw94u0wzc67qa03cmgl92qwqap09p8xa
Shares: 10000000000000000000/1Height: 215307
```

> Please notice that the share amount is also correspond to iris-atto, 1iris=10^18 iris-atto

* Re-delegate 

Once a delegator has delegated his own IRIS to certain validator, he/she could change the destination of delegation at anytime. If the transaction is executed, the 
delegation will be placed at the other's pool after the specified period of the system parameter `unbonding_time`. 
 
you should run the following command:
```$xslt
iriscli stake redelegate --addr-validator-dest=<addr-validator-dest>  --addr-validator-source=<addr-validator> --address-delegator=<address-delegator>  --chain-id=irishub  --from=name --gas=50000 --fee=0.3iris --shares-percent=1.0 
```

The example output is the following:
```$xslt
Delegation
Delegator: iaa1je9qyff4qate4e0kthum0p8v7q7z8lr7phygv8
Validator: iaa1kepndxvjr6gnc8tjcnelp9hqz8jdcs8m5dcr88
Shares: 10000000000000000000/1Height: 215459
```

* Unbond Delegation

Once a delegator has delegated his own IRIS to certain validator, he/she could withdraw the delegation at anytime. If the transaction is executed, the 
delegation will become liquid after after the specified period of the system parameter `unbonding_time`.  

To start, you should run the following command:
```$xslt
iriscli stake unbond  begin  --address-validator=<address-validator> --address-delegator=<address-delegator>  --chain-id=irishub  --from=name --gas=50000 --fee=0.3iris --shares-percent=1.0 
```

You could check that the balance of delegator has increased.
