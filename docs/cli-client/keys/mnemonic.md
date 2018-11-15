# iriscli keys mnemonic

## Description

Create a bip39 mnemonic, sometimes called a seed phrase, by reading from the system entropy. To pass your own entropy, use --unsafe-entropy

## Usage

```
iriscli keys mnemonic <name> [flags]
```

## Flags

| Name, shorthand  | Default   | Description                                                                   | Required |
| ---------------- | --------- | ----------------------------------------------------------------------------- | -------- |
| --help, -h       |           | help for add                                                                  |          |
| --unsafe-entropy |           | Prompt the user to supply their own entropy, instead of relying on the system |          |

## Examples

### Create a bip39 mnemonic

```shell
iriscli keys mnemonic abc
```

After this, you'll get a bip39 mnemonic.

```txt
police possible oval milk network indicate usual blossom spring wasp taste canal announce purpose rib mind river pet brown web response sting remain airport
```