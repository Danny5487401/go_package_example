<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Neo4j之Cypher语句](#neo4j%E4%B9%8Bcypher%E8%AF%AD%E5%8F%A5)
    - [CREATE 语句：](#create-%E8%AF%AD%E5%8F%A5)
    - [MATCH语句](#match%E8%AF%AD%E5%8F%A5)
    - [RETURN 语句：](#return-%E8%AF%AD%E5%8F%A5)
    - [关系创建](#%E5%85%B3%E7%B3%BB%E5%88%9B%E5%BB%BA)
      - [1. 在两个现有的节点之间创建无属性的关系](#1-%E5%9C%A8%E4%B8%A4%E4%B8%AA%E7%8E%B0%E6%9C%89%E7%9A%84%E8%8A%82%E7%82%B9%E4%B9%8B%E9%97%B4%E5%88%9B%E5%BB%BA%E6%97%A0%E5%B1%9E%E6%80%A7%E7%9A%84%E5%85%B3%E7%B3%BB)
      - [2. 在两个现有的节点之间创建有属性的关系](#2-%E5%9C%A8%E4%B8%A4%E4%B8%AA%E7%8E%B0%E6%9C%89%E7%9A%84%E8%8A%82%E7%82%B9%E4%B9%8B%E9%97%B4%E5%88%9B%E5%BB%BA%E6%9C%89%E5%B1%9E%E6%80%A7%E7%9A%84%E5%85%B3%E7%B3%BB)
      - [3. 在两个新节点之间创建无属性的关系](#3-%E5%9C%A8%E4%B8%A4%E4%B8%AA%E6%96%B0%E8%8A%82%E7%82%B9%E4%B9%8B%E9%97%B4%E5%88%9B%E5%BB%BA%E6%97%A0%E5%B1%9E%E6%80%A7%E7%9A%84%E5%85%B3%E7%B3%BB)
      - [4. 在两个新节点之间创建有属性的关系](#4-%E5%9C%A8%E4%B8%A4%E4%B8%AA%E6%96%B0%E8%8A%82%E7%82%B9%E4%B9%8B%E9%97%B4%E5%88%9B%E5%BB%BA%E6%9C%89%E5%B1%9E%E6%80%A7%E7%9A%84%E5%85%B3%E7%B3%BB)
    - [WHERE 子句](#where-%E5%AD%90%E5%8F%A5)
    - [REMOVE 子句](#remove-%E5%AD%90%E5%8F%A5)
    - [SET 子句](#set-%E5%AD%90%E5%8F%A5)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Neo4j之Cypher语句
### CREATE 语句：
- 创建没有属性的节点
- 使用属性创建节点
- 在没有属性的节点之间创建关系
- 使用属性创建节点之间的关系
- 为节点或关系创建单个或多个标签

创建没有属性的节点：
```
#格式：
CREATE(<node_name>:<label_name>)
#例子：
CREATE(stu:Student)

```

创建有属性的节点:（注：属性值得用‘’或“”标注）
```
#格式：
CREATE(
	<node-name>:<label-name>
	{
		<Property1-name>:<Property1-Value>，(每个属性之间用逗号分开）
			.......
		<Propertyn-name>:<Propertyn-Value>

	}
)
#例子：
CREATE(stu:Student{name:'danny',age:26})

```
### MATCH语句
- 从数据库获取有关节点和属性的数据
- 从数据库获取有关节点，关系和属性的数据

```cypher
#格式：
MATCH(
<node-name>:<label-name>
)
#例子：
MATCH(stu:Student)
```
Note: 如果我们单独使用它，那么我们将InvalidSyntax错误。
如果你观察到错误消息，它告诉我们，我们可以使用MATCH命令与RETURN子句或更新子句

### RETURN 语句：
- 检索节点的某些属性 
- 检索节点的所有属性 
- 检索节点和关联关系的某些属性 
- 检索节点和关联关系的所有属性

### 关系创建
- 在两个现有的节点之间创建无属性的关系
- 在两个现有的节点之间创建有属性的关系
- 在两个新节点之间创建无属性的关系
- 在两个新节点之间创建有属性的关系
- 在具有WHERE子句的两个退出节点之间创建/不使用属性的关系

基本语法
1. – 表示一个无指向的关系
2. –> 表示一个有指向的关系
```
CREATE(:)-[]->(:)

[<relationship_name>:<relationship_label>
	{
<Property1-name>:<Property1-Value>
	}
] """"
eg:
[role]#只有名字
[role:ACTION_IN]#标签名字都有
[:ACTION_IN]#省略了名字，只有标签
[role:ACTION_IN{roles:'Advanced'}]#有名字，有标签，有属性等

```
note: 由于关系不能单独创建，所以这里需要与MATCH,CREATE等一起使用，即引出了模式这一说法

- （）表示头节点和尾节点
- [ ]表示关系

完整写法
```
CREATE(
    node_name:node_lable
    {Property-name:Property-Value})-
    [relationship_name:relationship_label]->
    (node_name:node_label
    {Property-name:Property-Value}
    )
eg:
CREATE(
	tea:Teacher:Actor{name: "素人"} )
	-[teach:ACTED_IN{roles: ["star"] } ]->
	(stu:Student{name: "Kurry"} 
	)
```
#### 1. 在两个现有的节点之间创建无属性的关系
```shell
MATCH (e:Customer),(cc:CreditCard) 
CREATE (e)-[r:DO_SHOPPING_WITH ]->(cc) 
```
这里关系名称为“r”,关系标签为“DO_SHOPPING_WITH”,e和Customer分别是客户节点的节点名称和节点标签名称。cc和CreditCard分别是CreditCard节点的节点名和节点标签名。

#### 2. 在两个现有的节点之间创建有属性的关系
```shell
MATCH (cust:Customer),(cc:CreditCard) 
CREATE (cust)-[r:DO_SHOPPING_WITH{shopdate:"12/12/2014",price:55000}]->(cc) 
RETURN r
```
#### 3. 在两个新节点之间创建无属性的关系
```shell
CREATE (fb1:FaceBookProfile1)-[like:LIKES]->(fb2:FaceBookProfile2)
```

#### 4. 在两个新节点之间创建有属性的关系
```shell
CREATE (video1:YoutubeVideo1{title:"Action Movie1",updated_by:"Abc",uploaded_date:"10/10/2010"})
-[movie:ACTION_MOVIES{rating:1}]->
(video2:YoutubeVideo2{title:"Action Movie2",updated_by:"Xyz",uploaded_date:"12/12/2012"})

```

### WHERE 子句
```shell
WHERE <condition>
WHERE <condition> <boolean-operator> <condition>

```
note: 
- CQL中的布尔运算符:or,and,not,xor
- 比较运算符:==,<,>,=
condition语法
```shell
<property-name> <comparison-operator> <value>

prepared:
    CREATE (dept:Dept { deptno:10,dname:"Accounting",location:"Hyderabad" })
    CREATE (emp1:Employee{id:123,name:"Lokesh",sal:35000,deptno:10})
    CREATE (emp2:Employee{id:124,name:"sqwe",sal:36000,deptno:20})
    CREATE (emp3:Employee{id:125,name:"qwer",sal:37000,deptno:30})
    CREATE (emp4:Employee{id:126,name:"asdf",sal:38000,deptno:40})

eg:
    MATCH (emp:Employee) 
    WHERE emp.name = 'asdf' OR emp.name = 'sqwe'
    RETURN emp
```

### REMOVE 子句
REMOVE子句主要用于删除节点或关系的属性，这点跟DELETE有些差别。
DELETE是删除节点或关系，但不删除属性。

DELETE和REMOVE命令之间的主要区别 

- DELETE操作用于删除节点和关联关系。
- REMOVE操作用于删除标签和属性。

```
// CREATE (book:Book {id:122,title:"Neo4j Tutorial",pages:340,price:250}) 
// match (b:Book) return b;
match (b:Book{id:122}) remove b.price return b;

// 删除很多个属性,用逗号（，）分隔开
MATCH (book { id:122 })
REMOVE book.price，book.id，book.pages
RETURN book

MATCH (m:Movie) 
REMOVE m:Picture
```

### SET 子句
- 向现有节点或关系添加新属性
- 添加或更新属性值

```shell
MATCH (dc:DebitCard)
SET dc.atm_pin = 3456
RETURN dc
```