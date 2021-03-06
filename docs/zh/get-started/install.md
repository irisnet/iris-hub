---
order: 2
---

# 安装
## 安装`go`

:::tip
编译安装 IRIShub 软件依赖 **Go 1.15+**。
:::

按照[官方文档](https://golang.org/doc/install) 安装`go`。

别忘记设置您的`$GOPATH`，`$GOBIN`和`$PATH`环境变量，例如：

```bash
mkdir -p $HOME/go/bin
echo "export GOPATH=$HOME/go" >> ~/.bashrc
source ~/.bashrc
echo "export GOBIN=$GOPATH/bin" >> ~/.bashrc
source ~/.bashrc
echo "export PATH=$PATH:$GOBIN" >> ~/.bashrc
source ~/.bashrc
```

确认已成功安装`go`

```bash
go version
```

## 安装`iris`

正确配置`go`之后，您应该可以编译并运行`iris`了。

请确保您的服务器可以访问 google.com，因为我们的项目依赖于google提供的某些库（如果您无法访问`google.com`，也可以尝试添加代理：`export GOPROXY=https://goproxy.io`）

```bash
git clone https://github.com/irisnet/irishub
cd irishub
git checkout <version>
make install
```

如果环境变量配置无误，则通过运行以上命令即可完成`iris`的安装。现在检查您的`iris`版本是否正确：

```bash
iris version
```
