<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [expr-lang/expr](#expr-langexpr)
  - [使用](#%E4%BD%BF%E7%94%A8)
    - [内置函数](#%E5%86%85%E7%BD%AE%E5%87%BD%E6%95%B0)
    - [Predicate](#predicate)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [第三方使用- argo-rollout](#%E7%AC%AC%E4%B8%89%E6%96%B9%E4%BD%BF%E7%94%A8--argo-rollout)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# expr-lang/expr


Expr表达式引擎是一个针对Go语言设计的动态配置解决方案，它以简单的语法和强大的性能特性著称。
Expr表达式引擎的核心是安全、快速和直观，很适合用于处理诸如访问控制、数据过滤和资源管理等场景。在Go语言中应用Expr，可以极大地提升应用程序处理动态规则的能力。
不同于其他语言的解释器或脚本引擎，Expr采用了静态类型检查，并且生成字节码来执行，因此它能同时保证性能和安全性


## 使用


### 内置函数
```go
var Builtins = []*Function{
	{
		Name:      "all",
		Predicate: true,
		Types:     types(new(func([]any, func(any) bool) bool)),
	},
	{
		Name:      "none",
		Predicate: true,
		Types:     types(new(func([]any, func(any) bool) bool)),
	},
	// ...
}
```


array 相关
- 函数 all 可以用来检验集合中的元素是否全部满足给定的条件。它接受两个参数，第一个参数是集合，第二个参数是条件表达式。
```go
// 检查所有 tweets 的 Content 长度是否小于 240
code := `all(tweets, len(.Content) < 240)`
```

- any 函数用来检测集合中是否有任一元素满足条件。

```go
// 检查是否有任一 tweet 的 Content 长度大于 240
code := `any(tweets, len(.Content) > 240)`
```


- one 函数用于确认集合中只有一个元素满足条件。
```go
// 检查是否只有一个 tweet 包含了特定关键词
code := `one(tweets, contains(.Content, "关键词"))`
```


### Predicate 

```shell
filter(0..9, {# % 2 == 0})
```


## 源码分析

```shell
(⎈|test-ctx:istio-system)➜  expr@v1.16.9 tree -L 2 .        
.
├── LICENSE
├── README.md
├── SECURITY.md
├── ast
│   ├── dump.go
│   ├── node.go
│   ├── print.go
│   ├── print_test.go
│   ├── visitor.go
│   └── visitor_test.go
├── bench_test.go
├── builtin
│   ├── builtin.go
│   ├── builtin_test.go
│   ├── function.go
│   ├── lib.go
│   ├── utils.go
│   └── validation.go
├── checker
│   ├── checker.go
│   ├── checker_test.go
│   ├── info.go
│   ├── info_test.go
│   └── types.go
├── compiler
│   ├── compiler.go
│   └── compiler_test.go
├── conf
│   ├── config.go
│   └── types_table.go
├── docgen
│   ├── README.md
│   ├── docgen.go
│   ├── docgen_test.go
│   └── markdown.go
├── docs
│   ├── configuration.md
│   ├── environment.md
│   ├── functions.md
│   ├── getting-started.md
│   ├── language-definition.md
│   ├── patch.md
│   └── visitor.md
├── expr.go
├── expr_test.go
├── file
│   ├── error.go
│   ├── location.go
│   ├── source.go
│   └── source_test.go
├── go.mod
├── internal
│   ├── deref
│   ├── difflib
│   ├── spew
│   └── testify
├── optimizer
│   ├── const_expr.go
│   ├── filter_first.go
│   ├── filter_last.go
│   ├── filter_len.go
│   ├── filter_map.go
│   ├── fold.go
│   ├── in_array.go
│   ├── in_range.go
│   ├── optimizer.go
│   ├── optimizer_test.go
│   ├── predicate_combination.go
│   ├── sum_array.go
│   ├── sum_array_test.go
│   ├── sum_map.go
│   └── sum_map_test.go
├── parser
│   ├── lexer
│   ├── operator
│   ├── parser.go
│   ├── parser_test.go
│   └── utils
├── patcher
│   ├── operator_override.go
│   ├── value
│   ├── with_context.go
│   ├── with_context_test.go
│   ├── with_timezone.go
│   └── with_timezone_test.go
├── test
│   ├── coredns
│   ├── crowdsec
│   ├── deref
│   ├── fuzz
│   ├── gen
│   ├── interface_method
│   ├── mock
│   ├── operator
│   ├── patch
│   ├── pipes
│   ├── playground
│   └── time
├── testdata
│   ├── crash.txt
│   ├── crowdsec.json
│   └── examples.txt
└── vm
    ├── debug.go
    ├── debug_off.go
    ├── debug_test.go
    ├── func_types
    ├── func_types[generated].go
    ├── opcodes.go
    ├── program.go
    ├── program_test.go
    ├── runtime
    ├── utils.go
    ├── vm.go
    └── vm_test.go

37 directories, 78 files

```


```go
// 解析并编译成字节码
func Compile(input string, ops ...Option) (*vm.Program, error) {
	// 初始化配置：包括内置函数
	config := conf.CreateNew()
	for _, op := range ops {
		op(config)
	}
	// 关闭指定的内置函数
	for name := range config.Disabled {
		delete(config.Builtins, name)
	}
	config.Check()

	// 根据传入的内容生成解析后的 ats 树
	tree, err := checker.ParseCheck(input, config)
	if err != nil {
		return nil, err
	}

	if config.Optimize {
		err = optimizer.Optimize(&tree.Node, config)
		if err != nil {
			var fileError *file.Error
			if errors.As(err, &fileError) {
				return nil, fileError.Bind(tree.Source)
			}
			return nil, err
		}
	}

	program, err := compiler.Compile(tree, config)
	if err != nil {
		return nil, err
	}

	return program, nil
}
```




## 第三方使用- argo-rollout


```go
// https://github.com/argoproj/argo-rollouts/blob/ff3471a2fc3ccb90dbb1f370d7e399ff3064043a/utils/evaluate/evaluate.go
func EvalCondition(resultValue interface{}, condition string) (bool, error) {
	var err error

	env := map[string]interface{}{
		"result":  valueFromPointer(resultValue),
		"asInt":   asInt,
		"asFloat": asFloat,
		"isNaN":   math.IsNaN,
		"isInf":   isInf,
		"isNil":   isNilFunc(resultValue),
		"default": defaultFunc(resultValue),
	}

	unwrapFileErr := func(e error) error {
		if fileErr, ok := err.(*file.Error); ok {
			e = errors.New(fileErr.Message)
		}
		return e
	}

	program, err := expr.Compile(condition, expr.Env(env))
	if err != nil {
		return false, unwrapFileErr(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, unwrapFileErr(err)
	}

	switch val := output.(type) {
	case bool:
		return val, nil
	default:
		return false, fmt.Errorf("expected bool, but got %T", val)
	}
}
```



## 参考

- 官方文档：https://expr-lang.org/docs/language-definition