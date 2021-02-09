# IRIShub 从 v0.16 迁移到 v1.0.0

## 1. 导出 genesis 文件

停止 irishub v0.16 守护程序，并使用 `irishub v0.16.4（已修复导出中的错误）` 在升级高度导出主网状态并指定 `--for-zero-height`。

```bash
iris export --home [v0.16_node_home] --height [upgrade-height] --for-zero-height
```

## 2. 迁移 genesis 文件

使用 `irishub v1.0.0` 迁移导出的 genesis.json 文件

```bash
iris migrate genesis.json --chain-id irishub-1 > genesis_v1.0.0.json
```

校验 sha256sum 哈希值是否正确

```bash
sha256sum genesis_v1.0.0.json
```

## 3. 初始化新节点

使用 `irishub v1.0.0` 初始化新的节点

```bash
iris init [moniker] --home [v1.0.0_node_home]
```

## 4. 迁移私钥文件

使用 `irishub v1.0.0` 迁移私钥文件

```bash
go run migrate/scripts/privValUpgrade.go [v0.16_node_home]/config/priv_validator.json [v1.0.0_node_home]/config/priv_validator_key.json [v1.0.0_node_home]/data/priv_validator_state.json
```

## 5. 迁移节点密钥文件

使用 `irishub v1.0.0` 迁移节点密钥文件

```bash
cp [v0.16_node_home]/config/node_key.json [v1.0.0_node_home]/config/node_key.json
```

## 6. 复制迁移过的 genensis 文件到新的节点目录

复制 genesis_v1.0.0.json 到新节点的配置文件目录

```bash
cp genesis_v1.0.0.json [v1.0.0_node_home]/config/genesis.json
```

## 7. 配置新的节点

在 `[v1.0.0_node_home]/config/app.toml` 中配置 `minimum-gas-prices`

```toml

# The minimum gas prices a validator is willing to accept for processing a
# transaction. A transaction's fees must meet the minimum of any denomination
# specified in this config (e.g. 0.25token1;0.0001token2).
minimum-gas-prices = "0.2uiris"

```

复制 `[v0.16_node_home]/config/config.toml` 中的 `persistent_peers` 到 `[v1.0.0_node_home]/config/config.toml` 中

```toml

# Comma separated list of nodes to keep persistent connections to
persistent_peers = ""

```

参考原节点配置文件 `[v0.16_node_home]/config/config.toml` 配置新的节点

## 8. 启动新的节点

使用 irishub v1.0.0 启动新的节点

```bash
iris unsafe-reset-all --home [v1.0.0_node_home]
iris start --home [v1.0.0_node_home]
```
