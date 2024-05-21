
# Masterminds/squirrel



比ORM还是要复杂一些,却又比裸写SQL好一些(可维护性好一些，不容易出 SQL注入问题)


## 源码分析

squirrel的四大结构体

- SelectBuilder
- UpdateBuilder
- InsertBuilder
- DeleteBuilde

这里拿 Select 进行说明

```go
// github.com/!masterminds/squirrel@v1.5.4/statement.go

// StatementBuilderType is the type of StatementBuilder.
type StatementBuilderType builder.Builder

// StatementBuilder is a parent builder for other builders, e.g. SelectBuilder.
var StatementBuilder = StatementBuilderType(builder.EmptyBuilder).PlaceholderFormat(Question)

func Select(columns ...string) SelectBuilder {
	return StatementBuilder.Select(columns...)
}

```