<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go.uber.org/dig](#gouberorgdig)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [注入依赖](#%E6%B3%A8%E5%85%A5%E4%BE%9D%E8%B5%96)
    - [Invoke 开始注入](#invoke-%E5%BC%80%E5%A7%8B%E6%B3%A8%E5%85%A5)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go.uber.org/dig


通过反射的方式来实现运行时的依赖注入.


## 源码分析

### 初始化

生成一个容器，这个容器是由一个个scope组成的有向无环图
```go
//  go.uber.org/dig@v1.19.0/container.go
func New(opts ...Option) *Container {
	s := newScope()
	c := &Container{scope: s} // 树的根节点

	for _, opt := range opts {
		opt.applyOption(c)
	}
	return c
}


func newScope() *Scope {
	s := &Scope{
		providers:       make(map[key][]*constructorNode),
		decorators:      make(map[key]*decoratorNode),
		values:          make(map[key]reflect.Value),
		decoratedValues: make(map[key]reflect.Value),
		groups:          make(map[key][]reflect.Value),
		decoratedGroups: make(map[key]reflect.Value),
		invokerFn:       defaultInvoker,
		rand:            rand.New(rand.NewSource(time.Now().UnixNano())),
		clockSrc:        digclock.System,
	}
	s.gh = newGraphHolder(s)
	return s
}
```


### 注入依赖

```go
func (c *Container) Provide(constructor interface{}, opts ...ProvideOption) error {
	return c.scope.Provide(constructor, opts...)
}

func (s *Scope) Provide(constructor interface{}, opts ...ProvideOption) error {
    // ...

	if err := s.provide(constructor, options); err != nil {
        // ...
	}
	return nil
}



func (s *Scope) provide(ctor interface{}, opts provideOptions) (err error) {
	// If Export option is provided to the constructor, this should be injected to the
	// root-level Scope (Container) to allow it to propagate to all other Scopes.
	origScope := s
	if opts.Exported {
		s = s.rootScope()
	}

	// For all scopes affected by this change,
	// take a snapshot of the current graph state before
	// we start making changes to it as we may need to
	// undo them upon encountering errors.
	allScopes := s.appendSubscopes(nil)

	defer func(allSc []*Scope) {
		if err != nil {
			for _, sc := range allSc {
				sc.gh.Rollback()
			}
		}
	}(allScopes)

	for _, sc := range allScopes {
		sc.gh.Snapshot()
	}

	n, err := newConstructorNode(
		ctor,
		s,
		origScope,
		constructorOptions{
			ResultName:     opts.Name,
			ResultGroup:    opts.Group,
			ResultAs:       opts.As,
			Location:       opts.Location,
			Callback:       opts.Callback,
			BeforeCallback: opts.BeforeCallback,
		},
	)
	if err != nil {
		return err
	}

	keys, err := s.findAndValidateResults(n.ResultList())
	if err != nil {
		return err
	}

	ctype := reflect.TypeOf(ctor)
	if len(keys) == 0 {
		return newErrInvalidInput(
			fmt.Sprintf("%v must provide at least one non-error type", ctype), nil)
	}

	oldProviders := make(map[key][]*constructorNode)
	for k := range keys {
		// Cache old providers before running cycle detection.
		oldProviders[k] = s.providers[k]
		s.providers[k] = append(s.providers[k], n)
	}

	for _, s := range allScopes {
		s.isVerifiedAcyclic = false
		if s.deferAcyclicVerification {
			continue
		}
		if ok, cycle := graph.IsAcyclic(s.gh); !ok {
			// When a cycle is detected, recover the old providers to reset
			// the providers map back to what it was before this node was
			// introduced.
			for k, ops := range oldProviders {
				s.providers[k] = ops
			}

			return newErrInvalidInput("this function introduces a cycle", s.cycleDetectedError(cycle))
		}
		s.isVerifiedAcyclic = true
	}

	s.nodes = append(s.nodes, n)

	// Record introspection info for caller if Info option is specified
	if info := opts.Info; info != nil {
		params := n.ParamList().DotParam()
		results := n.ResultList().DotResult()

		info.ID = (ID)(n.id)
		info.Inputs = make([]*Input, len(params))
		info.Outputs = make([]*Output, len(results))

		for i, param := range params {
			info.Inputs[i] = &Input{
				t:        param.Type,
				optional: param.Optional,
				name:     param.Name,
				group:    param.Group,
			}
		}

		for i, res := range results {
			info.Outputs[i] = &Output{
				t:     res.Type,
				name:  res.Name,
				group: res.Group,
			}
		}
	}
	return nil
}

```

初始化 Node

```go
func newConstructorNode(ctor interface{}, s *Scope, origS *Scope, opts constructorOptions) (*constructorNode, error) {
	cval := reflect.ValueOf(ctor)
	ctype := cval.Type()
	cptr := cval.Pointer()

	// 获取函数的参数列表
	params, err := newParamList(ctype, s)
	if err != nil {
		return nil, err
	}

	results, err := newResultList(
		ctype,
		resultOptions{
			Name:  opts.ResultName,
			Group: opts.ResultGroup,
			As:    opts.ResultAs,
		},
	)
	if err != nil {
		return nil, err
	}

	location := opts.Location
	if location == nil {
		location = digreflect.InspectFunc(ctor)
	}

	n := &constructorNode{
		ctor:           ctor, // 函数
		ctype:          ctype, //函数反射信息
		location:       location, // 函数的信息
		id:             dot.CtorID(cptr),
		paramList:      params, // 函数参数信息
		resultList:     results,  //函数返回信息
		orders:         make(map[*Scope]int),
		s:              s,
		origS:          origS,
		callback:       opts.Callback,
		beforeCallback: opts.BeforeCallback,
	}
	s.newGraphNode(n, n.orders)
	return n, nil
}
```


参数类型

- paramSingle好理解，注入函数的一般形参比如int、string、struct、slice都属于paramSingle

- paramGroupedSlice组类型
```go
StudentList []*Student `group:"stu"`
```

- paramObject 嵌入dig.In的结构体类型.paramObject可以包含 paramSingle和paramGroupedSlice类型

```go
type DBInfo struct {
    dig.In
    PrimaryDSN   *DSN `name:"primary"`
    SecondaryDSN *DSN `name:"secondary"`
}
```

### Invoke 开始注入
```go
func (c *Container) Invoke(function interface{}, opts ...InvokeOption) error {
	return c.scope.Invoke(function, opts...)
}

func (s *Scope) Invoke(function interface{}, opts ...InvokeOption) (err error) {
	ftype := reflect.TypeOf(function)
	if ftype == nil {
		return newErrInvalidInput("can't invoke an untyped nil", nil)
	}
	if ftype.Kind() != reflect.Func {
		return newErrInvalidInput(
			fmt.Sprintf("can't invoke non-function %v (type %v)", function, ftype), nil)
	}

	// 解析参数
	pl, err := newParamList(ftype, s)
	if err != nil {
		return err
	}

	// 通过参数解析所需要的依赖
	if err := shallowCheckDependencies(s, pl); err != nil {
		return errMissingDependencies{
			Func:   digreflect.InspectFunc(function),
			Reason: err,
		}
	}

	// 判断是否有环
	if !s.isVerifiedAcyclic {
		if ok, cycle := graph.IsAcyclic(s.gh); !ok {
			return newErrInvalidInput("cycle detected in dependency graph", s.cycleDetectedError(cycle))
		}
		s.isVerifiedAcyclic = true
	}

	args, err := pl.BuildList(s)
	if err != nil {
		return errArgumentsFailed{
			Func:   digreflect.InspectFunc(function),
			Reason: err,
		}
	}
	if s.recoverFromPanics {
		defer func() {
			if p := recover(); p != nil {
				err = PanicError{
					fn:    digreflect.InspectFunc(function),
					Panic: p,
				}
			}
		}()
	}

	var options invokeOptions
	for _, o := range opts {
		o.applyInvokeOption(&options)
	}

	// Record info for the invoke if requested
	if info := options.Info; info != nil {
		params := pl.DotParam()
		info.Inputs = make([]*Input, len(params))
		for i, p := range params {
			info.Inputs[i] = &Input{
				t:        p.Type,
				optional: p.Optional,
				name:     p.Name,
				group:    p.Group,
			}
		}

	}

	returned := s.invokerFn(reflect.ValueOf(function), args)
	if len(returned) == 0 {
		return nil
	}
	if last := returned[len(returned)-1]; isError(last.Type()) {
		if err, _ := last.Interface().(error); err != nil {
			return err
		}
	}

	return nil
}

```

## 参考

- [分解uber依赖注入库dig-使用篇](https://www.cnblogs.com/li-peng/p/14708132.html)
- [分解uber依赖注入库dig-源码分析](https://www.cnblogs.com/li-peng/p/14738098.html)
