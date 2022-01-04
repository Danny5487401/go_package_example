# MakeFile

## 规则
```makefile
target: dependencies
    system command(s)
```
- target 通常是程序要生成的目标文件的名字. 但也可以是一个动作的名字.

- dependencies 是依赖, 通常是文件, 完成 target 所需要的输入.

- system command(s) 是完成 target 所需要运行的指令, 即 shell 命令.


一条语句一行, 使用单个 tab 缩进.使用 make 命令可以运行各种 target. 如果不带 target 参数,
第一个 target 会被作为默认目标.

Makefile 不是为了编译, 也不再引用任何文件, 仅仅是为了整合多个命令, 比写脚本方便多.
这个时候涉及到一个叫做伪目标的指令 .PHONY. 

.PHONY 后面跟着的多个 target 都不是要生成的文件的名字,
而是指代一个动作, 一个行为. 比如 test 指运行测试, clean 清理文件等

Note:windows 下没有 make 命令, 所以 Makefile 也就无法使用.

## 使用
### 1. 变量
```makefile
APP=myapp

build: clean
	go build -o ${APP} main.go

clean:
	rm -rf ${APP}
```

### 2. 递归的目标
当前工程目录结构
```css
~/project

├── main.go
├── Makefile
└── mymodule/
      ├── main.go
      └── Makefile
```

文件根目录下还有一个文件夹 mymodule，它可能是一个单独的模块，也需要打包构建，并且定义有自己的 Makefile :
```makefile
# ~/project/mymodule/Makefile

APP=module

build:
	go build -o ${APP} main.go
```

现在当你处于项目的根目录时，如何去执行 mymodule 子目录下定义的 Makefile 呢？

使用 cd 命令也可以，不过我们有其它的方式去解决这个问题：使用 -C 标志和特定的 ${MAKE} 变量
```makefile
APP=myapp

.PHONY: build
build: clean
	go build -o ${APP} main.go

.PHONY: clean
clean:
	rm -rf ${APP}


.PHONY: build-mymodule
build-mymodule:
	${MAKE} -C mymodule build
```
当你执行 make build-mymodule 时，其将会自动切换到 mymodule 目录，并且执行 mymodule 目录下的 Makefile 中定义的 build 指令

### 3.shell 输出作为变量
```makefile
V=$(shell go version)

gv:
	echo ${V}
```

### 4. 判断语句
假设我们的指令依赖于环境变量 ENV ，我们可以使用一个前置条件去检查是否忘了输入 ENV 
```makefile
.PHONY: run
run: check-env
	echo ${ENV}

check-env:
ifndef ENV
    $(error ENV not set, allowed values - `staging` or `production`)
endif
```

这里当我们执行 make run 时，因为有前置条件 check-env 会先执行前置条件中的内容，指令内容是一个判断语句，判断 ENV 是否未定义，如果未定义，则会抛出一个错误，错误提示就是 error 后面的内容

### 5.帮助提示
添加 help 帮助提示
```makefile
.PHONY: build
## build: build the application
build: clean
    @echo "Building..."
    @go build -o ${APP} main.go

.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
    @echo "Cleaning"
    @rm -rf ${APP}

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init \
		&& go mod tidy \
		&& go mod vendor

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'
```

这样当你执行 make help 时，就是打印如下的提示内容：
```css
Usage:

  build   build the application
  run     runs go run main.go
  clean   cleans the binary
  setup   setup go modules
  help    prints this help message
```