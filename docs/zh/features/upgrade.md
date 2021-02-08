# 软件升级

## 概念

### 计划

升级模块定义一个 `Plan`（`计划`） 类型，该计划调度一个在线升级过程。`计划` 可以安排在一个特定的区块高度或时间，但不能同时指定两者。一旦对一个（冻结的）候选发布版本以及一个合适的处理器达成一致，则一个 `计划` 被创建，其 `Name` 对应一个特定的处理器。通常，`计划` 通过一个治理提议过程被创建，当投票通过时，该计划将被调度。一个计划的 `Info` 可以包含关于这次升级的各种元数据，典型的，要包含一些上链的特定于应用的升级信息，诸如验证人能自动升级的 `Git` 提交。

#### 边车进程

如果一个运行应用程序二进制的运营者也运行一个边车进程来辅助自动下载和升级二进制，`Info` 允许这个进程是无摩擦的。即，升级模块实现 [cosmovisor 可升级的二进制规范](https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor#upgradeable-binary-specification) 指定的规范，并且 `cosmovisor` 能可选地被用于为节点运营者完全自动化升级过程。通过用必要的信息填充 `Info` 字段，二进制能够被自动下载。参考[这里](https://github.com/cosmos/cosmos-sdk/tree/master/cosmovisor#auto-download)。

```go
type Plan struct {
 Name string
 Time time.Time
 Height int64
 Info string
 UpgradedClientState *types.Any
}
```

### 处理器

`升级` 模块便于从主版本 X 升级到主版本 Y。为完成这个过程，节点运营者必须首先将当前的二进制升级到一个新二进制，该二进制有一个与新版本 Y 相对应的 `处理器`。假定这个版本已经过大部分社区成员的充分测试和批准。这个 `处理器` 定义了在新二进制 Y 成功运行链之前需要完成的状态迁移。当然，这个 `处理器` 是特定于应用而不是在模块基础上定义的。通过应用中的 `Keeper#SetUpgradeHandler` 完成 `处理器` 的注册。

```go
type UpgradeHandler func(Context, Plan)
```

在每个 `BeginBlock` 执行期间，`升级` 模块检查是否存在一个应该执行的 `计划`（被调度在 `BeginBlock` 运行时的区块高度或时间）。如果存在，则执行对应的 `处理器`。如果这个计划期待被执行但没有注册相应的处理器，或者二进制升级过早时，节点将优雅地 `panic` 并退出。

### 存储加载器

升级模块也有助于将存储迁移作为升级的一部分。`存储加载器` 执行在新二进制成功运行链之前需要完成的迁移。`存储加载器` 也是特定于应用而不是基于模块定义的。通过应用中的 `app#SetStoreLoader` 完成 `存储加载器` 的注册。

```go
func UpgradeStoreLoader (upgradeHeight int64, storeUpgrades *store.StoreUpgrades) baseapp.StoreLoader
```

如果存在一个计划的升级并且已到达升级高度，在 panic 之前旧的二进制将写 `升级信息（UpgradeInfo）` 到磁盘。

```go
type UpgradeInfo struct {
    Name    string
    Height  int64
}
```

这个信息对于确保 `存储升级` 在正确的区块高度顺利执行以及确保升级符合预期至关重要。它消除了新二进制每次重启时多次执行 `存储升级` 的机会。而且，如果在相同的高度存在多个升级计划，`Name` 将确保这些 `存储升级` 仅仅发生在计划的升级处理器中。

目前在升级过程中，涉及到的状态迁移方式主要支持以下三种：`Renamed`、`Deleted`、`Added`

#### Renamed

用户可以在升级过程中指定将oldKey(前缀)下的所有数据迁移到newKey(前缀)下存储。

#### Deleted

用户可以在升级过程中删除指定key(前缀)下的所有数据。

#### Added

用户可以在升级过程中以指定key为前缀申请一块新的存储区域。

### 提议

通常，一个 `计划` 经由 `软件升级提议（SoftwareUpgradeProposal）` 通过治理的方式被提出并提交，这个提议规定了标准的治理过程。如果这个提议被通过，目标为一个特定 `处理器` 的 `计划` 将被持久化并调度。通过在一个新的提议中更新计划时间（`Plan.Time`），升级可以被推迟或者加速。

```go
type SoftwareUpgradeProposal struct {
    Title       string
    Description string
    Plan        Plan
}
```

#### 治理过程

当一个升级提议被接受时，升级过程分为如下两个步骤。

##### 停止网络共识

软件升级提议被接受后，系统将在指定高度的`BeginBlock`阶段执行升级前的准备，包括下载升级计划、暂停网络共识。

###### 下载升级计划

为了能够顺利升级软件，必须在停止网络共识之前，先记录升级所需要的信息：`计划名称`、`升级高度`。
  
- `计划名称`：网络重启时，需要根据计划名称路由到对应的`UpgradeHandler`以及`UpgradeStoreLoader`。
- `升级高度`：网络重启时，检查是否需要执行网络升级计划。

###### 暂停网络共识

软件升级提议被接受后，系统将在指定高度的`BeginBlock`阶段优雅的暂停网络共识。

##### 重新启动新软件

用户替换软件为指定版本并重新启动网络，系统将检测是否包含`计划名称`所指定的`处理器`，如果包含，系统首先执行`处理器`程序，然后开始网络共识，如果不包含，系统报错并退出。

#### 取消升级提议

升级提议可以被取消。存在一个 `取消软件升级 （CancelSoftwareUpgrade）`的提议类型，当该类型提议投票通过时，将移除当前正在进行的升级计划。当然，这需要在升级计划执行之前被投票通过并执行。

如果当前升级计划已经被执行，但是该升级计划存在问题，那么此时提出`取消软件升级`提议是无效的(因为网络已经停止共识)。这时还有另一方案来弥补这一过失，就是在重新启动网络时，使用`--unsafe-skip-upgrades`参数跳过指定的升级高度(并不是真的跳过该高度，而是跳过软件升级`处理器`)。当然这要求参与共识的2/3验证人都执行同样的操作，否则同样无法达成网络共识。

## 升级流程

### 提交升级提案

执行软件升级流程的第一步是由治理模块发起一个软件升级提案，该提案详细说明了升级时间以及升级内容，具体见上面[概念](#概念)。发起提案的命令行示例如下：

```bash
iris tx gov submit-proposal software-upgrade bifrost-rc2 \
  --deposit 1000iris \
  --upgrade-time 2021-02-09T13:00:00Z \
  --title "mainnet software upgrade" \
  --upgrade-info "Commit: 0ef5dd0b4d140a4788f05fc1a0bd409b3c6a0492. After the proposal is approved, please use the commit hash to build and restart your node." \
  --description "Upgrade the mainnet software version from v1.0.0-rc0 to v1.0.0-rc2."
  --from=node0 --chain-id=test --fees=1iris -b block -y
```

### 为提案抵押、投票

软件升级提案和其他普通提案的执行流程基本一致，都需要验证人、委托人为该提案发表意见，具体信息请参考[治理模块](./governance.md)。为提案抵押的命令行示例如下：

```bash
iris tx gov deposit 1 1000iris --from=node0 --chain-id=test --fees=1iris -b block -y
```

一旦抵押金额达到最小抵押金额，提案将进入投票期，验证人或者委托人需要对该提案发起投票，发起投票的命令行示例如下：

```bash
iris tx gov vote 1 yes  --from=node0 --chain-id=test --fees=1iris -b block -y
```

当软件升级提案被通过后，升级模块会创建一项升级计划，在指定高度或者时间使所有节点停止网络共识，等待新的软件重启网络。

### 重启网络

当升级提案通过，节点会停止出块，用户需要根据在第一步[提交升级提案](#提交升级提案)中指明的新版本信息下载源码，编译新软件，具体参考[安装](./../get-started/install.md)。新软件安装完成后，使用新的版本重启节点，节点会执行和计划名称相应的升级逻辑。一旦全网超过2/3的权重使用新的版本重启网络，区块链网络将重新达成新的共识，继续产生新区块。
