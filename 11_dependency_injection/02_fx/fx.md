<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [go.uber.org/fx](#gouberorgfx)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [生命周期管理](#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F%E7%AE%A1%E7%90%86)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# go.uber.org/fx

Fx 是 Uber 开发并开源的Go语言模块组合框架，它提供了一种模块化、插拔、可组合的方式来构建Go应用。


## 源码分析


```go
// go.uber.org/fx@v1.24.0/app.go

func New(opts ...Option) *App {
	logger := fxlog.DefaultLogger(os.Stderr)  //获取日志实例

	app := &App{
		clock:        fxclock.System,
		startTimeout: DefaultTimeout,
		stopTimeout:  DefaultTimeout,
		receivers:    newSignalReceivers(),
	}
	app.root = &module{
		app: app,
		// We start with a logger that writes to stderr. One of the
		// following three things can change this:
		//
		// - fx.Logger was provided to change the output stream
		// - fx.WithLogger was provided to change the logger
		//   implementation
		// - Both, fx.Logger and fx.WithLogger were provided
		//
		// The first two cases are straightforward: we use what the
		// user gave us. For the last case, however, we need to fall
		// back to what was provided to fx.Logger if fx.WithLogger
		// fails.
		log:   logger,
		trace: []string{fxreflect.CallerStack(1, 2)[0].String()},
	}

	for _, opt := range opts {
		opt.apply(app.root)
	}

	// There are a few levels of wrapping on the lifecycle here. To quickly
	// cover them:
	//
	// - lifecycleWrapper ensures that we don't unintentionally expose the
	//   Start and Stop methods of the internal lifecycle.Lifecycle type
	// - lifecycleWrapper also adapts the internal lifecycle.Hook type into
	//   the public fx.Hook type.
	// - appLogger ensures that the lifecycle always logs events to the
	//   "current" logger associated with the fx.App.
	app.lifecycle = &lifecycleWrapper{
		lifecycle.New(appLogger{app}, app.clock),
	}

	containerOptions := []dig.Option{
		dig.DeferAcyclicVerification(),
		dig.DryRun(app.validate),
	}

	if app.recoverFromPanics {
		containerOptions = append(containerOptions, dig.RecoverFromPanics())
	}

	// 创建dig container
	app.container = dig.New(containerOptions...)
	app.root.build(app, app.container)

	// Provide Fx types first to increase the chance a custom logger
	// can be successfully built in the face of unrelated DI failure.
	// E.g., for a custom logger that relies on the Lifecycle type.
	frames := fxreflect.CallerStack(0, 0) // include New in the stack for default Provides
	app.root.provide(provide{
		Target: func() Lifecycle { return app.lifecycle },
		Stack:  frames,
	})
	app.root.provide(provide{Target: app.shutdowner, Stack: frames})
	app.root.provide(provide{Target: app.dotGraph, Stack: frames})
	app.root.provideAll()

	// Run decorators before executing any Invokes
	// (including the ones inside installAllEventLoggers).
	app.err = multierr.Append(app.err, app.root.decorateAll())

	// If you are thinking about returning here after provides: do not (just yet)!
	// If a custom logger was being used, we're still buffering messages.
	// We'll want to flush them to the logger.

	// custom app logger will be initialized by the root module.
	app.root.installAllEventLoggers()

	// This error might have come from the provide loop above. We've
	// already flushed to the custom logger, so we can return.
	if app.err != nil {
		return app
	}

	if err := app.root.invokeAll(); err != nil {
        // 错误处理
	}

	return app
}
```
## 生命周期管理



## modules 模块

## 参考


- https://uber-go.github.io/fx/get-started/
- [深入解析go依赖注入库go.uber.org/fx](https://zhuanlan.zhihu.com/p/418299054)
- [使用uber-go的fx进行依赖注入](https://czyt.tech/post/using-uber-go-fx-for-go-dependency-injection/)