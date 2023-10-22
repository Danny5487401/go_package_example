<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [设计思路](#%E8%AE%BE%E8%AE%A1%E6%80%9D%E8%B7%AF)
  - [lua前置知识](#lua%E5%89%8D%E7%BD%AE%E7%9F%A5%E8%AF%86)
  - [ui界面](#ui%E7%95%8C%E9%9D%A2)
  - [客户端-->生产者](#%E5%AE%A2%E6%88%B7%E7%AB%AF--%E7%94%9F%E4%BA%A7%E8%80%85)
  - [server -->消费者](#server---%E6%B6%88%E8%B4%B9%E8%80%85)
    - [状态转移](#%E7%8A%B6%E6%80%81%E8%BD%AC%E7%A7%BB)
    - [处理](#%E5%A4%84%E7%90%86)
      - [ProcessTask处理函数查找：type与方法的匹配-->参考go官方的http server的路由匹配](#processtask%E5%A4%84%E7%90%86%E5%87%BD%E6%95%B0%E6%9F%A5%E6%89%BEtype%E4%B8%8E%E6%96%B9%E6%B3%95%E7%9A%84%E5%8C%B9%E9%85%8D--%E5%8F%82%E8%80%83go%E5%AE%98%E6%96%B9%E7%9A%84http-server%E7%9A%84%E8%B7%AF%E7%94%B1%E5%8C%B9%E9%85%8D)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# 源码分析

## 设计思路
- 延迟队列设计思路：zset的分值为时间
- 消息有状态：活跃，计划中，重试，已完成等，状态迁移使用list，如果状态是已经完成的key需要删除
- 多个消息可以放到多个通道，防止拥堵，queue使用set
- 数据内容存储：hash


## lua前置知识
1. 私有变量，比如当前变量只在其对应的方法中起作用，就需要在声明前加上 " local"关键字
2. 关系运算符：不等于为  ~= , 特殊符号：  .. 这两点表示两个字符串相加
3. if语句
```lua
if(5 < 10)
then            --用then和end代替以往的{}进行包裹
    代码体
end 
```
4. table 是 Lua 语言中的一种“数据/代码结构,Lua 语言中的数组其实就是 table 类型
```lua

myArray = {'住址',1234, true}
 
table.getn(myArray)   --可通过该方法获得该数组长度

for key, value in ipairs(myArray) do        --如果是数组结构，用 ipairs 方法；如果是键值对结构，用 pairs 方法.
      print(key, value)
end
```

## ui界面
```shell
docker run --rm     --name asynqmon     -p 8080:8080     hibiken/asynqmon --redis-addr "172.17.0.1:6379"
```

## 客户端-->生产者
- set 存放所有的队列
```shell
127.0.0.1:6379> smembers  "asynq:queues"
1) "low"
2) "default"
```
```go
	if err := r.client.SAdd(ctx, base.AllQueues, msg.Queue).Err(); err != nil {
		return errors.E(op, errors.Unknown, &errors.RedisCommandError{Command: "sadd", Err: err})
	}
```

定时的lua脚本
```go
var scheduleCmd = redis.NewScript(`
if redis.call("EXISTS", KEYS[1]) == 1 then
	return 0
end
redis.call("HSET", KEYS[1],
           "msg", ARGV[1],
           "state", "scheduled")
redis.call("ZADD", KEYS[2], ARGV[2], ARGV[3])
return 1
`)
```
```go
// Schedule adds the task to the scheduled set to be processed in the future.
func (r *RDB) Schedule(ctx context.Context, msg *base.TaskMessage, processAt time.Time) error {
	var op errors.Op = "rdb.Schedule"
	// 使用proto转换数据
	encoded, err := base.EncodeMessage(msg)
	if err != nil {
		return errors.E(op, errors.Unknown, fmt.Sprintf("cannot encode message: %v", err))
	}
	// set保存queue名字
	if err := r.client.SAdd(ctx, base.AllQueues, msg.Queue).Err(); err != nil {
		return errors.E(op, errors.Unknown, &errors.RedisCommandError{Command: "sadd", Err: err})
	}
	
	// lua脚本中的Key和value
	keys := []string{
		base.TaskKey(msg.Queue, msg.ID),
		base.ScheduledKey(msg.Queue),
	}
	argv := []interface{}{
		encoded,
		processAt.Unix(),
		msg.ID,
	}
	n, err := r.runScriptWithErrorCode(ctx, op, scheduleCmd, keys, argv...)
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.E(op, errors.AlreadyExists, errors.ErrTaskIdConflict)
	}
	return nil
}
```
- hash存储: key是队列名包含唯一uid， 字段msg,状态state,
```shell
127.0.0.1:6379> type "asynq:{default}:t:42cf52dd-865c-45d9-b63c-7a5a4d1e64b0"
hash
127.0.0.1:6379> hgetall "asynq:{default}:t:42cf52dd-865c-45d9-b63c-7a5a4d1e64b0"
1) "msg"
2) "\n\remail:deliver\x12-{\"UserID\":42,\"TemplateID\":\"some:template:id\"}\x1a$42cf52dd-865c-45d9-b63c-7a5a4d1e64b0\"\adefault(\x19@\x88\x0e"
3) "pending_since"
4) "1652793345345129000"
5) "state"
6) "pending"
```

- zset存储：key "asynq:{low}:scheduled" ， 成员id , 分值 时间s
```shell
127.0.0.1:6379> type "asynq:{low}:scheduled"
zset
127.0.0.1:6379> zrevrange "asynq:{low}:scheduled" 0 -1 withscores
1) "2777c58c-510e-4628-9064-49f740cacdcf"
2) "1652793405"
```
注意：  队列名字 + 状态


## server -->消费者

删除时的lua脚本
```go
var deleteExpiredCompletedTasksCmd = redis.NewScript(`
local ids = redis.call("ZRANGEBYSCORE", KEYS[1], "-inf", ARGV[1], "LIMIT", 0, tonumber(ARGV[3]))
for _, id in ipairs(ids) do
	redis.call("DEL", ARGV[2] .. id)
	redis.call("ZREM", KEYS[1], id)
end
return table.getn(ids)`)
```
```go
func (r *RDB) deleteExpiredCompletedTasks(qname string, batchSize int) (int64, error) {
	var op errors.Op = "rdb.DeleteExpiredCompletedTasks"
	keys := []string{base.CompletedKey(qname)}
	argv := []interface{}{
		r.clock.Now().Unix(),
		base.TaskKeyPrefix(qname),
		batchSize,
	}
	res, err := deleteExpiredCompletedTasksCmd.Run(context.Background(), r.client, keys, argv...).Result()
	if err != nil {
		return 0, errors.E(op, errors.Internal, fmt.Sprintf("redis eval error: %v", err))
	}
	n, ok := res.(int64)
	if !ok {
		return 0, errors.E(op, errors.Internal, fmt.Sprintf("unexpected return value from Lua script: %v", res))
	}
	return n, nil
}
```
根据小的分数返回zset中数据是id，然后删除hash中的数据，然后zset中的数据





### 状态转移
```go
// KEYS[1] -> source queue (e.g. asynq:{<qname>:scheduled or asynq:{<qname>}:retry})
// KEYS[2] -> asynq:{<qname>}:pending
// ARGV[1] -> current unix time in seconds
// ARGV[2] -> task key prefix
// ARGV[3] -> current unix time in nsec
// ARGV[4] -> group key prefix
// Note: Script moves tasks up to 100 at a time to keep the runtime of script short.
var forwardCmd = redis.NewScript(`
local ids = redis.call("ZRANGEBYSCORE", KEYS[1], "-inf", ARGV[1], "LIMIT", 0, 100)
for _, id in ipairs(ids) do
	local taskKey = ARGV[2] .. id
	local group = redis.call("HGET", taskKey, "group")
	if group and group ~= '' then
	    redis.call("ZADD", ARGV[4] .. group, ARGV[1], id)
		redis.call("ZREM", KEYS[1], id)
		redis.call("HSET", taskKey,
				   "state", "aggregating")
	else
		redis.call("LPUSH", KEYS[2], id)
		redis.call("ZREM", KEYS[1], id)
		redis.call("HSET", taskKey,
				   "state", "pending",
				   "pending_since", ARGV[3])
	end
end
return table.getn(ids)`)
```

主要看scheduled-->pending


### 处理
出队数据
````go
// Input:
// KEYS[1] -> asynq:{<qname>}:pending
// KEYS[2] -> asynq:{<qname>}:paused
// KEYS[3] -> asynq:{<qname>}:active
// KEYS[4] -> asynq:{<qname>}:lease
// --
// ARGV[1] -> initial lease expiration Unix time
// ARGV[2] -> task key prefix
//
// Output:
// Returns nil if no processable task is found in the given queue.
// Returns an encoded TaskMessage.
//
// Note: dequeueCmd checks whether a queue is paused first, before
// calling RPOPLPUSH to pop a task from the queue.
var dequeueCmd = redis.NewScript(`
if redis.call("EXISTS", KEYS[2]) == 0 then
	local id = redis.call("RPOPLPUSH", KEYS[1], KEYS[3])
	if id then
		local key = ARGV[2] .. id
		redis.call("HSET", key, "state", "active")
		redis.call("HDEL", key, "pending_since")
		redis.call("ZADD", KEYS[4], ARGV[1], id)
		return redis.call("HGET", key, "msg")
	end
end
return nil`)
````
```go
func (r *RDB) Dequeue(qnames ...string) (msg *base.TaskMessage, leaseExpirationTime time.Time, err error) {
	var op errors.Op = "rdb.Dequeue"
	for _, qname := range qnames {
		keys := []string{
			base.PendingKey(qname),
			base.PausedKey(qname),
			base.ActiveKey(qname),
			base.LeaseKey(qname),
		}
		leaseExpirationTime = r.clock.Now().Add(LeaseDuration)
		argv := []interface{}{
			leaseExpirationTime.Unix(),
			base.TaskKeyPrefix(qname),
		}
		res, err := dequeueCmd.Run(context.Background(), r.client, keys, argv...).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, time.Time{}, errors.E(op, errors.Unknown, fmt.Sprintf("redis eval error: %v", err))
		}
		encoded, err := cast.ToStringE(res)
		if err != nil {
			return nil, time.Time{}, errors.E(op, errors.Internal, fmt.Sprintf("cast error: unexpected return value from Lua script: %v", res))
		}
		// 还原数据
		if msg, err = base.DecodeMessage([]byte(encoded)); err != nil {
			return nil, time.Time{}, errors.E(op, errors.Internal, fmt.Sprintf("cannot decode message: %v", err))
		}
		return msg, leaseExpirationTime, nil
	}
	return nil, time.Time{}, errors.E(op, errors.NotFound, errors.ErrNoProcessableTask)
}
```


#### ProcessTask处理函数查找：type与方法的匹配-->参考go官方的http server的路由匹配

路由匹配-->type名字匹配
```go
func (mux *ServeMux) Handler(t *Task) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	h, pattern = mux.match(t.Type())
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	for i := len(mux.mws) - 1; i >= 0; i-- {
		h = mux.mws[i](h)
	}
	return h, pattern
}
```