<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [curd](#curd)
  - [数据类型](#%E6%95%B0%E6%8D%AE%E7%B1%BB%E5%9E%8B)
  - [create](#create)
  - [insert](#insert)
  - [update and delete](#update-and-delete)
  - [clickhouse-client 使用](#clickhouse-client-%E4%BD%BF%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# curd

默认情况下，CREATE、DROP、ALTER 和 RENAME 查询仅影响执行它们的当前服务器。在集群设置中，可以使用 ON CLUSTER 子句以分布式方式运行此类查询。


## 数据类型

- (Int) 和 (unsigned UInt):不区分 int,short, long

- bool: 实际是UInt8.

- Float32/Float64: 这会精度丢失,建议用 Decimal


- Decimal: 建议只使用 Decimal32/Decimal64, 因为Decimal128由计算机模拟

- String/FixedString : 没有 Varchar,blob,clob类型


- Enum: 实际存储 Int8 or Int16

- array

- 时间: Date( 1970-01-01),DateTime(2019-01-01 00:00:00 ),DateTime64(2019-01-01 03:00:00.123)

- Nullable:  It returns 1 if the corresponding value is NULL and 0 otherwise. 建议业务上,字符串用空替代Null, 整型用 -1 表示 Null
```shell
CREATE TABLE nullable (`n` Nullable(UInt32)) ENGINE = MergeTree ORDER BY tuple();

INSERT INTO nullable VALUES (1) (NULL) (2) (NULL);

SELECT n.null FROM nullable;
┌─n.null─┐
│      0 │
│      1 │
│      0 │
│      1 │
└────────┘
```

## create

```clickhouse
-- 数据库创建
CREATE DATABASE [IF NOT EXISTS] db_name [ON CLUSTER cluster] [ENGINE = engine(...)] [COMMENT 'Comment']

-- 创建本地表
CREATE TABLE [IF NOT EXISTS] [db.]table_name ON CLUSTER cluster
(
  name1 [type1] [DEFAULT|MATERIALIZED|ALIAS expr1],
  name2 [type2] [DEFAULT|MATERIALIZED|ALIAS expr2],
  ...
  INDEX index_name1 expr1 TYPE type1(...) GRANULARITY value1,
  INDEX index_name2 expr2 TYPE type2(...) GRANULARITY value2
  ) ENGINE = engine_name()
  [PARTITION BY expr]
  [ORDER BY expr]
  [PRIMARY KEY expr]
  [SAMPLE BY expr]
  [SETTINGS name=value, ...];

```

- MATERIALIZED：物化列表达式，表示该列不能被INSERT，是被计算出来的； 在INSERT语句中，不需要写入该列；在SELECT * 查询语句结果集不包含该列；需要指定列表来查询（虚拟列）
- ALIAS ：别名列。这样的列不会存储在表中。 它的值不能够通过INSERT写入，同时SELECT查询使用星号时，这些列也不会被用来替换星号。 但是它们可以用于SELECT中，在这种情况下，在查询分析中别名将被替换

物化列与别名列的区别： 物化列是会保存数据，查询的时候不需要计算，而别名列不会保存数据，查询的时候需要计算，查询时候返回表达式的计算结果

```clickhouse
-- 创建一个本地表
CREATE TABLE ontime_local ON CLUSTER default -- 表名为 ontime_local
(
    Year UInt16,
    Quarter UInt8,
    Month UInt8,
    DayofMonth UInt8,
    DayOfWeek UInt8,
    FlightDate Date,
    FlightNum String,
    Div5WheelsOff String,
    Div5TailNum String
)ENGINE = ReplicatedMergeTree(--表引擎用ReplicatedMergeTree，开启数据副本的合并树表引擎）
    '/clickhouse/tables/ontime_local/{shard}', -- 指定存储路径
    '{replica}')           
 PARTITION BY toYYYYMM(FlightDate)  -- 指定分区键，按FlightDate日期转年+月维度，每月做一个分区
 PRIMARY KEY (intHash32(FlightDate)) -- 指定主键，FlightDate日期转hash值
 ORDER BY (intHash32(FlightDate),FlightNum) -- 指定排序键，包含两列：FlightDate日期转hash值、FlightNunm字符串。
 SAMPLE BY intHash32(FlightDate)  -- 抽样表达式，采用FlightDate日期转hash值
SETTINGS index_granularity= 8192 ;  -- 指定index_granularity指数，每个分区再次划分的数量

```

```clickhouse
-- 基于本地表创建一个分布式表
CREATE TABLE  [db.]table_name  ON CLUSTER default
 AS db.local_table_name
ENGINE = Distributed(<cluster>, <database>, <shard table> [, sharding_key])

```
- sharding_key：分片表达式。可以是一个字段，例如user_id（integer类型），通过对余数值进行取余分片；也可以是一个表达式，例如rand()，通过rand()函数返回值/shards总权重分片；为了分片更均匀，可以加上hash函数，如intHash64(user_id)
```clickhouse
CREATE TABLE ontime_distributed ON CLUSTER default   -- 指定分布式表的表名，所在集群
 AS db_name.ontime_local                             -- 指定对应的 本地表的表名
ENGINE = Distributed(default, db_name, ontime_local, rand());  -- 指定表引擎为Distributed（固定）

```



## insert


插入最佳
- 建议一次插入一千行数据,理想是1000~10000行.
- 插入可以重试,因为是幂等的.对于 MergeTree engine ,会自动去重.



## update and delete 
clickhouse 没有直接支持 update and delete
注意：

- 索引列不支持更新、删除
- 分布式表不支持更新、删除
```clickhouse
ALTER TABLE [<database>.]<table> UPDATE <column> = <expression> WHERE <filter_expr>
```




```clickhouse
ALTER TABLE [<database>.]<table> DELETE WHERE <filter_expr>
```

轻量删除
- 默认异步.Set mutations_sync equal to 1 to wait for one replica to process the statement, and set mutations_sync to 2 to wait for all replicas
- 只适合MergeTree
```clickhouse
DELETE FROM [db.]table [ON CLUSTER cluster] [WHERE expr]
```



## clickhouse-client 使用

```shell
I have no name!@my-clickhouse-shard0-0:/$ clickhouse-client --help
Main options:
  --help                                  print usage summary, combine with --verbose to display all options
  --verbose                               print query and other debugging info
  -V [ --version ]                        print version information and exit
  --version-clean                         print version in machine-readable format and exit
  -C [ --config-file ] arg                config-file path
  -q [ --query ] arg                      Query. Can be specified multiple times (--query "SELECT 1" --query "SELECT 2") or once with multiple comma-separated queries (--query "SELECT
                                          1; SELECT 2;"). In the latter case, INSERT queries with non-VALUE format must be separated by empty lines.
  --queries-file arg                      file path with queries to execute; multiple files can be specified (--queries-file file1 file2...)
  -n [ --multiquery ]                     Obsolete, does nothing
  -m [ --multiline ]                      If specified, allow multiline queries (do not send the query on Enter)
  -d [ --database ] arg                   database
  --query_kind arg (=initial_query)       One of initial_query/secondary_query/no_query
  --query_id arg                          query_id
  --history_file arg                      path to history file
  --stage arg (=complete)                 Request query processing up to specified stage: complete,fetch_columns,with_mergeable_state,with_mergeable_state_after_aggregation,with_merge
                                          able_state_after_aggregation_and_limit
  --progress [=arg(=tty)] (=default)      Print progress of queries execution - to TTY: tty|on|1|true|yes; to STDERR non-interactive mode: err; OFF: off|0|false|no; DEFAULT -
                                          interactive to TTY, non-interactive is off
  -A [ --disable_suggestion ]             Disable loading suggestion data. Note that suggestion data is loaded asynchronously through a second connection to ClickHouse server. Also it
                                          is reasonable to disable suggestion if you want to paste a query with TAB characters. Shorthand option -A is for those who get used to mysql
                                          client.
  --wait_for_suggestions_to_load          Load suggestion data synchonously.
  -t [ --time ]                           print query execution time to stderr in non-interactive mode (for benchmarks)
  --memory-usage [=arg(=default)] (=none) print memory usage to stderr in non-interactive mode (for benchmarks). Values: 'none', 'default', 'readable'
  --echo                                  in batch mode, print query before execution
  --log-level arg                         log level
  --server_logs_file arg                  put server logs into specified file
  --suggestion_limit arg (=10000)         Suggestion limit for how many databases, tables and columns to fetch.
  -f [ --format ] arg                     default output format (and input format for clickhouse-local)
  --output-format arg                     default output format (this option has preference over --format)
  -E [ --vertical ]                       vertical output format, same as --format=Vertical or FORMAT Vertical or \G at end of command
  --highlight arg (=1)                    enable or disable basic syntax highlight in interactive command line
  --ignore-error                          do not stop processing when an error occurs
  --stacktrace                            print stack traces of exceptions
  --hardware-utilization                  print hardware utilization information in progress bar
  --print-profile-events                  Printing ProfileEvents packets
  --profile-events-delay-ms arg (=0)      Delay between printing `ProfileEvents` packets (-1 - print only totals, 0 - print every single packet)
  --processed-rows                        print the number of locally processed rows
  --interactive                           Process queries-file or --query query and start interactive mode
  --pager arg                             Pipe all output into this command (less or similar)
  --max_memory_usage_in_client arg        Set memory limit in client/local server
  --client_logs_file arg                  Path to a file for writing client logs. Currently we only have fatal logs (when the client crashes)
  -c [ --config ] arg                     config-file path (another shorthand)
  --connection arg                        connection to use (from the client config), by default connection name is hostname
  -s [ --secure ]                         Use TLS connection
  --no-secure                             Don't use TLS connection
  -u [ --user ] arg (=default)            user
  --password arg                          password
  --ask-password                          ask-password
  --ssh-key-file arg                      File containing the SSH private key for authenticate with the server.
  --ssh-key-passphrase arg                Passphrase for the SSH private key specified by --ssh-key-file.
  --quota_key arg                         A string to differentiate quotas when the user have keyed quotas configured on server
  --jwt arg                               Use JWT for authentication
  --max_client_network_bandwidth arg      the maximum speed of data exchange over the network for the client in bytes per second.
  --compression arg                       enable or disable compression (enabled by default for remote communication and disabled for localhost communication).
  --query-fuzzer-runs arg (=0)            After executing every SELECT query, do random mutations in it and run again specified number of times. This is used for testing to discover
                                          unexpected corner cases.
  --create-query-fuzzer-runs arg (=0)
  --interleave-queries-file arg           file path with queries to execute before every file from 'queries-file'; multiple files can be specified (--queries-file file1 file2...);
                                          this is needed to enable more aggressive fuzzing of newly added tests (see 'query-fuzzer-runs' option)
  --opentelemetry-traceparent arg         OpenTelemetry traceparent header as described by W3C Trace Context recommendation
  --opentelemetry-tracestate arg          OpenTelemetry tracestate header as described by W3C Trace Context recommendation
  --no-warnings                           disable warnings when client connects to server
  --fake-drop                             Ignore all DROP queries, should be used only for testing
  --accept-invalid-certificate            Ignore certificate verification errors, equal to config parameters openSSL.client.invalidCertificateHandler.name=AcceptCertificateHandler and
                                          openSSL.client.verificationMode=none

External tables options:
  --file arg                   data file or - for stdin
  --name arg (=_data)          name of the table
  --format arg (=TabSeparated) data format
  --structure arg              structure
  --types arg                  types

Hosts and ports options:
  -h [ --host ] arg (=localhost) Server hostname. Multiple hosts can be passed via multiple argumentsExample of usage: '--host host1 --host host2 --port port2 --host host3 ...'Each
                                 '--port port' will be attached to the last seen host that doesn't have a port yet,if there is no such host, the port will be attached to the next
                                 first host or to default host.
  --port arg                     server ports

In addition, --param_name=value can be specified for substitution of parameters for parametrized queries.
```

## 参考
- https://clickhouse.com/docs/zh/sql-reference/statements