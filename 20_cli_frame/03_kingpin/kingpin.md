<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/alecthomas/kingpin](#githubcomalecthomaskingpin)
  - [特点](#%E7%89%B9%E7%82%B9)
  - [源码](#%E6%BA%90%E7%A0%81)
    - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [解析](#%E8%A7%A3%E6%9E%90)
  - [第三方使用](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/alecthomas/kingpin

当前项目主要靠其他人提 PR.



## 特点

如果用flag实现的话应该是下面这种使用方法：
```shell
./cli --method A
./cli --method B
./cli --method C
```

每次都需要输入"–method"，然而用kingpin库实现的话就可以达到下面这种效果：
```shell
./cli A
./cli B
./cli C
```

特点包括：
- 命令行参数定义：Kingpin 支持以链式调用的方式定义命令行参数，可以指定其名称、类型、默认值、描述等。
- 子命令和嵌套命令：Kingpin 允许创建多级的命令行结构，使得命令行工具可以有更好的组织结构。
- 命令行参数校验：Kingpin 提供了校验器，可以验证命令行参数的有效性。你可以指定自定义的校验函数，以确保用户输入的数据符合要求。
- 参数类型转换：支持常见的基本类型（int、string、float 等），以及复杂类型（如 []string、map[string]int）。
- 命令行错误处理：内置错误处理，帮助生成有用的错误信息。



## 源码

Application 结构体

```go
type Application struct {
    // 最关键的信息在 cmd mixin 
	cmdMixin 
	initialized bool // 是否已经 init ,即调用 app.parseContext

	Name string
	Help string

    // ...
	
	// 帮助信息模版
	usageTemplate  string
	usageFuncs     template.FuncMap
	validator      ApplicationValidator
	terminate      func(status int) // See Terminate()
	noInterspersed bool             // can flags be interspersed with args (or must they come first)
	defaultEnvars  bool
	completion     bool

    // ..
}

type cmdMixin struct {
	*flagGroup
	*argGroup
	*cmdGroup
	actionMixin
}
```

### 初始化

```go
// github.com/alecthomas/kingpin/v2@v2.4.0/app.go
func New(name, help string) *Application {
	a := &Application{
		Name:          name, // 应用名
		Help:          help,
		errorWriter:   os.Stderr, // Left for backwards compatibility purposes.
		usageWriter:   os.Stderr,
		usageTemplate: DefaultUsageTemplate,
		terminate:     os.Exit,
	}
	// flag 相关
	a.flagGroup = newFlagGroup()
	// arg 相关
	a.argGroup = newArgGroup()
	// cmd 相关
	a.cmdGroup = newCmdGroup(a)
	a.HelpFlag = a.Flag("help", "Show context-sensitive help (also try --help-long and --help-man).")
	a.HelpFlag.Bool()
    // ...

	return a
}
```

flag 添加

```go
func (f *flagGroup) Flag(name, help string) *FlagClause {
	flag := newFlag(name, help)
	f.long[name] = flag
	f.flagOrder = append(f.flagOrder, flag)
	return flag
}


func newFlag(name, help string) *FlagClause {
	f := &FlagClause{
		name: name,
		help: help,
	}
	return f
}

```

cmd 添加 

```go
func (a *Application) Command(name, help string) *CmdClause {
	return a.addCommand(name, help)
}


func (c *cmdGroup) addCommand(name, help string) *CmdClause {
	cmd := newCommand(c.app, name, help)
	c.commands[name] = cmd
	c.commandOrder = append(c.commandOrder, cmd)
	return cmd
}

```


arg 信息

```go
type ArgClause struct {
	actionMixin
	parserMixin
	completionsMixin
	envarMixin // 环境变量
	name          string
	help          string
	defaultValues []string //默认值
	placeholder   string
	hidden        bool
	required      bool //是否必须
}

```

arg 添加
```go
func (a *argGroup) Arg(name, help string) *ArgClause {
	arg := newArg(name, help)
	a.args = append(a.args, arg)
	return arg
}

func newArg(name, help string) *ArgClause {
	a := &ArgClause{
		name: name,
		help: help,
	}
	return a
}
```


### 解析

```go
func (a *Application) Parse(args []string) (command string, err error) {
    // 解析 
	context, parseErr := a.ParseContext(args)
	selected := []string{}
	var setValuesErr error

	if context == nil {
		// Since we do not throw error immediately, there could be a case
		// where a context returns nil. Protect against that.
		return "", parseErr
	}

	if err = a.setDefaults(context); err != nil {
		return "", err
	}

	selected, setValuesErr = a.setValues(context)

	if err = a.applyPreActions(context, !a.completion); err != nil {
		return "", err
	}

	if a.completion {
		a.generateBashCompletion(context)
		a.terminate(0)
	} else {
		if parseErr != nil {
			return "", parseErr
		}

		a.maybeHelp(context)
		if !context.EOL() {
			return "", fmt.Errorf("unexpected argument '%s'", context.Peek())
		}

		if setValuesErr != nil {
			return "", setValuesErr
		}

		command, err = a.execute(context, selected)
		if err == ErrCommandNotSpecified {
			a.writeUsage(context, nil)
		}
	}
	return command, err
}
```


token 解析

```go
func (p *ParseContext) Next() *Token {
	if len(p.peek) > 0 {
		return p.pop()
	}

	// End of tokens.
	if len(p.args) == 0 {
		return &Token{Index: p.argi, Type: TokenEOL}
	}

	if p.argi > 0 && p.argi <= len(p.rawArgs) && p.rawArgs[p.argi-1] == "--" {
		// If the previous argument was a --, from now on only arguments are parsed.
		p.argsOnly = true
	}
	arg := p.args[0]
	p.next()

	if p.argsOnly {
		return &Token{p.argi, TokenArg, arg}
	}

	if arg == "--" {
		return p.Next()
	}

	if strings.HasPrefix(arg, "--") {
		parts := strings.SplitN(arg[2:], "=", 2)
		token := &Token{p.argi, TokenLong, parts[0]}
		if len(parts) == 2 {
			p.Push(&Token{p.argi, TokenArg, parts[1]})
		}
		return token
	}

	if strings.HasPrefix(arg, "-") {
		if len(arg) == 1 {
			return &Token{Index: p.argi, Type: TokenArg}
		}
		shortRune, size := utf8.DecodeRuneInString(arg[1:])
		short := string(shortRune)
		flag, ok := p.flags.short[short]
		// Not a known short flag, we'll just return it anyway.
		if !ok {
		} else if fb, ok := flag.value.(boolFlag); ok && fb.IsBoolFlag() {
			// Bool short flag.
		} else {
			// Short flag with combined argument: -fARG
			token := &Token{p.argi, TokenShort, short}
			if len(arg) > size+1 {
				p.Push(&Token{p.argi, TokenArg, arg[size+1:]})
			}
			return token
		}

		if len(arg) > size+1 {
			p.args = append([]string{"-" + arg[size+1:]}, p.args...)
		}
		return &Token{p.argi, TokenShort, short}
	} else if EnableFileExpansion && strings.HasPrefix(arg, "@") {
		expanded, err := ExpandArgsFromFile(arg[1:])
		if err != nil {
			return &Token{p.argi, TokenError, err.Error()}
		}
		if len(p.args) == 0 {
			p.args = expanded
		} else {
			p.args = append(expanded, p.args...)
		}
		return p.Next()
	}

	return &Token{p.argi, TokenArg, arg}
}
```

## 第三方使用

- teleport 堡垒机
- node_exporter 

```go
// https://github.com/prometheus/node_exporter/blob/v1.8.2/node_exporter.go
func main() {
	// 使用 kingpin 定义了众多 flag
	var (
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		disableExporterMetrics = kingpin.Flag(
			"web.disable-exporter-metrics",
			"Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).",
		).Bool()
		maxRequests = kingpin.Flag(
			"web.max-requests",
			"Maximum number of parallel scrape requests. Use 0 to disable.",
		).Default("40").Int()
		disableDefaultCollectors = kingpin.Flag(
			"collector.disable-defaults",
			"Set all collectors to disabled by default.",
		).Default("false").Bool()
		maxProcs = kingpin.Flag(
			"runtime.gomaxprocs", "The target number of CPUs Go will run on (GOMAXPROCS)",
		).Envar("GOMAXPROCS").Default("1").Int()
		toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":9100")
	)

	promslogConfig := &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, promslogConfig)
	kingpin.Version(version.Print("node_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promslog.New(promslogConfig)

	if *disableDefaultCollectors {
		collector.DisableDefaultCollectors()
	}
	logger.Info("Starting node_exporter", "version", version.Info())
	logger.Info("Build context", "build_context", version.BuildContext())
	if user, err := user.Current(); err == nil && user.Uid == "0" {
		logger.Warn("Node Exporter is running as root user. This exporter is designed to run as unprivileged user, root is not required.")
	}
    // ..

	server := &http.Server{}
	if err := web.ListenAndServe(server, toolkitFlags, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
```



默认 parse
```go
// github.com/alecthomas/kingpin/v2@v2.4.0/global.go

func Parse() string {
	selected := MustParse(CommandLine.Parse(os.Args[1:]))
	if selected == "" && CommandLine.cmdGroup.have() {
		Usage()
		CommandLine.terminate(0)
	}
	return selected
}


var (
	// CommandLine 默认使用命令行应用名为 app name
	CommandLine = New(filepath.Base(os.Args[0]), "")
)
```




## 参考
