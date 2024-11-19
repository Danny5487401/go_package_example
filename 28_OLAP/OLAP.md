<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [数据仓库OLAP（Online Analytical Processing）](#%E6%95%B0%E6%8D%AE%E4%BB%93%E5%BA%93olaponline-analytical-processing)
  - [多维度分析案例](#%E5%A4%9A%E7%BB%B4%E5%BA%A6%E5%88%86%E6%9E%90%E6%A1%88%E4%BE%8B)
    - [销售明细表](#%E9%94%80%E5%94%AE%E6%98%8E%E7%BB%86%E8%A1%A8)
  - [OLAP分类](#olap%E5%88%86%E7%B1%BB)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 数据仓库 OLAP（Online Analytical Processing）
OLAP 名为联机分析处理，又可以称之为多维分析处理，指的是通过多种不同的维度审视数据，进行深层次分析.
例如clickhouse，greenplum，Doris。

- Greenplum是基于数据库分布式架构的开源大数据平台，采用无共享(no sharing)的MPP架构，具有良好的线性扩展能力，具有高效的并行运算、并行存储等特性。拥有独特的高效的ORCA优化器。兼容SQL语法。适合用于高效PB数据量级的存储、处理和实时分析能力。由于内核是基于PostgreSQL数据库，也支持涵盖OLTP型业务混合负载。同时数据节点和主节点都有自己备份节点。提供数据库的高可用性


|   比较内容   | OLTP（on-line transaction processing）联机事务处理 | OLAP（On-Line Analytical Processing）联机分析处理 |
|:--------:|:------------------------------------------:|:-----------------------------------------:|
|   操作特点   |             日常业务操作，尤其是包含大量前台操作             |            后台操作，例如统计报表，大批量数据加载            |
|   响应速度   |              优先级最高，要求相应速度非常高               |                要求速度高，吞吐量大                 |
|   吞吐量    |                     小                      |                     大                     |
|  并发访问量   |                    非常高                     |                    不高                     |
| 单笔事务消耗资源 |                     小                      |                     大                     |
| SQL 语句类型 |               插入和修改操作为主，DML                |              大量查询操作或批量DML操作               |
|     索引类型     |                    B*索引                    |           Bitmap、Bitmap Join 索引           |
|     索引量     |                     适量                     |                     多                     |
|     访问方式     |                    按索引访问                   |                   全表扫描                    |
|     连接方式     |                    Nested_loop                   |                   Hash Join                   |
|     BIND 变量     |                    	使用或强制使用	                 |                   不使用                    |

## 多维度分析案例

### 销售明细表
![](.OLAP_images/sale_detail.png)

1. 下钻：从高层次向低层次明细数据进行穿透。例如从 "省" 下钻到 "市"，从 "湖北省" 穿透到 "武汉" 和 "宜昌
![](.OLAP_images/sale_cubic1.png)

2. 和下钻相反，从低层次向高层次汇聚。例如从 "市" 汇聚到 "省"，将 "武汉" 和 "宜昌" 汇聚成 "湖北"
![](.OLAP_images/sale_cubic2.png)
3. 切片：观察立方体的一层，将一个或多个温度设为单个固定的值，然后观察剩余的维度，例如将商品维度固定为 "足球"。
![](.OLAP_images/sale_cubic3.png)
4. 切块：和切片类似，只是将单个固定值变成多个固定值。例如将商品维度固定为"足球"、"篮球" 和 "乒乓球"。
![](.OLAP_images/sale_cubic4.png)
5. 旋转：旋转立方体的一面，如果要将数据映射到一张二维表，那么就要进行旋转，等同于行列转换。
![](.OLAP_images/sale_cubic5.png)


## OLAP分类
1. 第一类架构称为 ROLAP（Relational OLAP，关系型 OLAP），顾名思义，它直接使用关系模型构建，数据模型常使用星型模型或者雪花模型，这是最先能够想到、也是最为直接的实现方法
2. 第二类架构称为 MOLAP（Multidimensional OLAP，多维型 OLAP），它的出现就是为了缓解 ROLAP 性能问题。
   - MOLAP 使用多维数组的形式存数据，其核心思想是借助预先聚合结果（说白了就是提前先算好，然后将结果保存起来），使用空间换取时间的形式从而提升查询性能。 
   也就是说，用更多的存储空间换得查询时间的减少，其具体的实现方式是依托立方体模型的概念。
   首先，对需要分析的数据进行建模，框定需要分析的维度字段；
   然后，通过预处理的形式，对各种维度进行组合并事先聚合；
   最后，将聚合结果以某种索引或者缓存的形式保存起来（通常只保留聚合后的结果，不存储明细数据），这样一来，在随后的查询过程中，可以直接利用结果返回数据。
3. 第三类架构称为 HOLAP（Hybrid OLAP，混合架构的OLAP），这种思路可以理解成 ROLAP 和 MOLAP 两者的组合