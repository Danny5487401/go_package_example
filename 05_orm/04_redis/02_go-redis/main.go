package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var redisdb *redis.Client
var wg sync.WaitGroup

func main() {
	wg.Add(1)
	go testRedisBase()
	wg.Wait()
}

func testRedisBase() {
	defer wg.Done()

	//连接服务器
	redisdb = redis.NewClient(&redis.Options{
		Addr:     "ali.danny.games:6379", // use default Addr
		Password: "",                     // no password set
		DB:       1,                      // use default DB
	})

	//心跳
	pong, err := redisdb.Ping().Result()
	log.Println(pong, err) // Output: PONG <nil>

	ExampleClient_String()
	//ExampleClient_List()
	//ExampleClient_Hash()
	//ExampleClient_Set()
	//ExampleClient_SortSet()
	//ExampleClient_HyperLogLog()
	//ExampleClient_CMD()
	//ExampleClient_Scan()
	//ExampleClient_Tx() // 事物pipeline
	//ExampleClient_Script()
	//ExampleClient_PubSub()

}

func ExampleClient_String() {
	log.Println("ExampleClient_String starts")
	defer log.Println("ExampleClient_String ends")

	//kv读写
	err := redisdb.Set("key", "value", 10*time.Second).Err()
	log.Println(err)

	// 判断key是否存在,不存在为0
	res, err := redisdb.Exists("set_key1").Result()
	log.Println("判断key是否存在", res, err)

	//获取过期时间
	tm, err := redisdb.TTL("key").Result()
	log.Println(tm)

	val, err := redisdb.Get("key_notExist").Result()
	if err == redis.Nil {
		log.Println("key不存在")
	}
	log.Println(val, err)

	val2, err := redisdb.Get("missing_key").Result()
	// redis.Nil 用于区分an empty string reply 和 a nil reply (key does not exist)键不存在:
	if err == redis.Nil {
		log.Println("missing_key does not exist")
	} else if err != nil {
		log.Println("missing_key", val2, err)
	}

	//不存在才设置 过期时间 nx ex
	value, err := redisdb.SetNX("counter", 0, 1*time.Second).Result()
	log.Println("setnx", value, err)

	//Incr
	result, err := redisdb.Incr("counter").Result()
	log.Println("Incr", result, err)
}

func ExampleClient_List() {
	log.Println("ExampleClient_List")
	defer log.Println("ExampleClient_List")

	//添加
	log.Println(redisdb.RPush("list_test", "message1").Err())
	log.Println(redisdb.RPush("list_test", "message2").Err())
	log.Println(redisdb.RPush("list_test", "message3").Err())
	log.Println(redisdb.RPush("list_test", "message4").Err())

	//设置
	log.Println(redisdb.LSet("list_test", 2, "message_set").Err())

	//remove
	ret, err := redisdb.LRem("list_test", 1, "message2").Result()
	log.Println(ret, err)

	rLen, err := redisdb.LLen("list_test").Result()
	log.Println(rLen, err)

	//遍历
	lists, err := redisdb.LRange("list_test", 0, rLen-1).Result()
	log.Println("LRange", lists, err)

	////pop没有时阻塞
	result, err := redisdb.BLPop(1*time.Second, "list_test").Result()
	log.Println("result:", result, err, len(result))
}

func ExampleClient_Hash() {
	log.Println("ExampleClient_Hash")
	defer log.Println("ExampleClient_Hash")

	datas := map[string]interface{}{
		"name": "danny",
		"sex":  1,
		"age":  28,
		"tel":  12345678,
	}

	//添加
	if err := redisdb.HMSet("hash_test", datas).Err(); err != nil {
		log.Fatal(err)
	}

	//获取
	rets, err := redisdb.HMGet("hash_test", "name", "sex").Result()
	log.Println("rets:", rets, err)

	//成员
	retAll, err := redisdb.HGetAll("hash_test").Result()
	log.Println("retAll", retAll, err)

	//存在
	bExist, err := redisdb.HExists("hash_test", "tel").Result()
	log.Println(bExist, err)
	// 只有在字段 field 不存在时，设置哈希表字段的值
	bRet, err := redisdb.HSetNX("hash_test", "id", 100).Result()
	log.Println(bRet, err)

	//删除
	log.Println(redisdb.HDel("hash_test", "age").Result())

	//为哈希表 key 中的指定字段的整数值加上增量 increment
	HRet, err := redisdb.HIncrBy("hash_test", "id", 10).Result()
	log.Println(HRet, err) // 返回的是增加后的结果
}

func ExampleClient_Set() {
	log.Println("ExampleClient_Set")
	defer log.Println("ExampleClient_Set")

	//第一次添加
	ret, err := redisdb.SAdd("set_test", "11", "22", "33", "44").Result()
	log.Println("第一次添加结果", ret, err) //
	//第一次添加
	ret2, err := redisdb.SAdd("set_test", "44").Result()
	log.Println("第二次添加结果", ret2, err) //返回的是成功添加的个数,如果是1，代表之前没有,不需要去exist判断

	//数量
	count, err := redisdb.SCard("set_test").Result()
	log.Println(count, err)

	//删除
	ret, err = redisdb.SRem("set_test", "11", "22").Result()
	log.Println(ret, err)

	//成员
	members, err := redisdb.SMembers("set_test").Result()
	log.Println(members, err)

	bret, err := redisdb.SIsMember("set_test", "100").Result()
	log.Println("是否是其成员", bret, err)

	redisdb.SAdd("set_a", "11", "22", "33", "44")
	redisdb.SAdd("set_b", "11", "22", "33", "55", "66", "77")
	//差集
	diff, err := redisdb.SDiff("set_a", "set_b").Result()
	log.Println(diff, err)

	//交集
	inter, err := redisdb.SInter("set_a", "set_b").Result()
	log.Println(inter, err)

	//并集
	union, err := redisdb.SUnion("set_a", "set_b").Result()
	log.Println(union, err)

	ret, err = redisdb.SDiffStore("set_diff", "set_a", "set_b").Result()
	log.Println(ret, err)

	rets, err := redisdb.SMembers("set_diff").Result()
	log.Println(rets, err)
}

func ExampleClient_SortSet() {
	log.Println("ExampleClient_SortSet")
	defer log.Println("ExampleClient_SortSet")

	addArgs := make([]redis.Z, 100)
	for i := 1; i < 100; i++ {
		addArgs = append(addArgs, redis.Z{Score: float64(i), Member: fmt.Sprintf("a_%d", i)})
	}
	//log.Println(addArgs)

	Shuffle := func(slice []redis.Z) {
		r := rand.New(rand.NewSource(time.Now().Unix()))
		for len(slice) > 0 {
			n := len(slice)
			randIndex := r.Intn(n)
			slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
			slice = slice[:n-1]
		}
	}

	//随机打乱
	Shuffle(addArgs)

	//添加
	ret, err := redisdb.ZAddNX("sortset_test", addArgs...).Result()
	log.Println("ZAddNX", ret, err)

	//获取指定成员score
	score, err := redisdb.ZScore("sortset_test", "a_10").Result()
	log.Println("ZScore", score, err)

	//获取制定成员的索引
	index, err := redisdb.ZRank("sortset_test", "a_50").Result()
	log.Println("ZRank", index, err)
	// ZCARD key 获取有序集合的成员数
	count, err := redisdb.ZCard("sortset_test").Result()
	log.Println("SCard", count, err)

	//返回有序集合指定区间内的成员
	rets, err := redisdb.ZRange("sortset_test", 10, 20).Result()
	log.Println(rets, err)

	//返回有序集合指定区间内的成员分数从高到低
	rets, err = redisdb.ZRevRange("sortset_test", 10, 20).Result()
	log.Println(rets, err)

	//指定分数区间的成员列表
	rets, err = redisdb.ZRangeByScore("sortset_test", redis.ZRangeBy{Min: "(30", Max: "(50", Offset: 1, Count: 10}).Result()
	log.Println(rets, err)
}

//用来做基数统计的算法，HyperLogLog 的优点是，在输入元素的数量或者体积非常非常大时，计算基数所需的空间总是固定 的、并且是很小的。
//每个 HyperLogLog 键只需要花费 12 KB 内存，就可以计算接近 2^64 个不同元素的基 数
func ExampleClient_HyperLogLog() {
	log.Println("ExampleClient_HyperLogLog")
	defer log.Println("ExampleClient_HyperLogLog")

	for i := 0; i < 10000; i++ {
		redisdb.PFAdd("pf_test_1", fmt.Sprintf("pfkey%d", i))
	}
	ret, err := redisdb.PFCount("pf_test_1").Result()
	log.Println(ret, err)

	for i := 0; i < 10000; i++ {
		redisdb.PFAdd("pf_test_2", fmt.Sprintf("pfkey%d", i))
	}
	ret, err = redisdb.PFCount("pf_test_2").Result()
	log.Println(ret, err)

	redisdb.PFMerge("pf_test", "pf_test_2", "pf_test_1")
	ret, err = redisdb.PFCount("pf_test").Result()
	log.Println(ret, err)
}

func ExampleClient_PubSub() {
	log.Println("ExampleClient_PubSub")
	defer log.Println("ExampleClient_PubSub")
	//发布订阅
	pubsub := redisdb.Subscribe("subkey")
	_, err := pubsub.Receive()
	if err != nil {
		log.Fatal("pubsub.Receive")
	}
	ch := pubsub.Channel()
	time.AfterFunc(1*time.Second, func() {
		log.Println("Publish")

		err = redisdb.Publish("subkey", "test publish 1").Err()
		if err != nil {
			log.Fatal("redisdb.Publish", err)
		}

		redisdb.Publish("subkey", "test publish 2")
	})
	for msg := range ch {
		log.Println("recv channel:", msg.Channel, msg.Pattern, msg.Payload)
	}
}

func ExampleClient_CMD() {
	log.Println("ExampleClient_CMD")
	defer log.Println("ExampleClient_CMD")

	//执行自定义redis命令
	Get := func(rdb *redis.Client, key string) *redis.StringCmd {
		cmd := redis.NewStringCmd("get", key)
		redisdb.Process(cmd)
		return cmd
	}

	v, err := Get(redisdb, "NewStringCmd").Result()
	log.Println("NewStringCmd", v, err)

	v, err = redisdb.Do("get", "redisdb.do").String()
	log.Println("redisdb.Do", v, err)
}

func ExampleClient_Scan() {
	log.Println("ExampleClient_Scan")
	defer log.Println("ExampleClient_Scan")

	//scan
	for i := 1; i < 1000; i++ {
		redisdb.Set(fmt.Sprintf("skey_%d", i), i, 0)
	}

	cusor := uint64(0)
	for {
		keys, retCusor, err := redisdb.Scan(cusor, "skey_*", int64(100)).Result()
		log.Println(keys, cusor, err)
		cusor = retCusor
		if cusor == 0 {
			break
		}
	}
}

// 事务pipeline
func ExampleClient_Tx() {
	pipe := redisdb.TxPipeline()
	incr := pipe.Incr("tx_pipeline_counter")
	boolCmd := pipe.Expire("tx_pipeline_counter", time.Hour)
	stringCmd := pipe.HGetAll("danny")

	// Execute
	//
	//     MULTI
	//     INCR pipeline_counter
	//     EXPIRE pipeline_counts 3600
	//     EXEC
	//
	// using one rdb-server roundtrip.
	cmders, err := pipe.Exec()
	fmt.Println(incr.Val(), boolCmd.Val(), stringCmd.Val(), err)
	for _, cmder := range cmders {
		switch cmd := cmder.(type) {
		case *redis.StringStringMapCmd:
			strMap, err := cmd.Result()
			if err != nil {
				fmt.Println("err", err)
			}
			fmt.Println("strMap", strMap)
		case *redis.IntCmd:
			intRet, err := cmd.Result()
			if err != nil {
				fmt.Println("err", err)
			}
			fmt.Println("IntRet", intRet)
		case *redis.BoolCmd:
			boolRet, err := cmd.Result()
			if err != nil {
				fmt.Println("err", err)
			}
			fmt.Println("boolRet", boolRet)

		}

	}
}

func ExampleClient_Script() {
	IncrByXX := redis.NewScript(`
        if redis.call("GET", KEYS[1]) ~= false then
            return redis.call("INCRBY", KEYS[1], ARGV[1])
        end
        return false
    `)

	n, err := IncrByXX.Run(redisdb, []string{"xx_counter"}, 2).Result()
	fmt.Println(n, err)

	err = redisdb.Set("xx_counter", "40", 0).Err()
	if err != nil {
		panic(err)
	}

	n, err = IncrByXX.Run(redisdb, []string{"xx_counter"}, 2).Result()
	fmt.Println(n, err)
}
