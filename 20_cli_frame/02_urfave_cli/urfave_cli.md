# urfave cli 


cli 目前已经开发到了 v2.0+。推荐使用最新的稳定版本。


对于 CLI 程序而言，我知道的最流行的框架有两个，分别是：

- urfave/cli：https://github.com/urfave/cli
- cobra：https://github.com/spf13/cobra
cobra 的功能会更强大完善。它的作者 Steve Francia（spf13）是 Google 里面 go 语言的 product lead，同时也是 gohugo、viper 等知名项目的作者。

但强大的同时，也意味着框架更大更复杂，在实现一些小规模的工具时，反而会觉得杀鸡牛刀

## 应用- buildkit
```go
// /Users/python/Downloads/git_download/buildkit/cmd/buildctl/main.go
func main() {
   ...
   app := cli.NewApp()
   app.Name = "buildctl"
   app.Usage = "build utility"
   app.Version = version.Version

   defaultAddress := os.Getenv("BUILDKIT_HOST")
   if defaultAddress == "" {
      defaultAddress = appdefaults.Address
   }

   app.Flags = []cli.Flag{
      ...
      cli.IntFlag{
         Name:  "timeout",
         Usage: "timeout backend connection after value seconds",
         Value: 5,
      },
   }

   app.Commands = []cli.Command{
      diskUsageCommand,
      pruneCommand,
      buildCommand,
      debugCommand,
      dialStdioCommand,
   }

   ...

   handleErr(debugEnabled, app.Run(os.Args))

```