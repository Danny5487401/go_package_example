package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/cenkalti/backoff/v4"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/esutil"
)

type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Published time.Time `json:"published"`
	Author    Author    `json:"author"`
}

type Author struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var (
	indexName  string
	numWorkers int
	flushBytes int
	numItems   int
)

func init() {
	flag.StringVar(&indexName, "index", "test-bulk-example", "索引名称")
	flag.IntVar(&numWorkers, "workers", runtime.NumCPU(), "工作进程数量")
	flag.IntVar(&flushBytes, "flush", 5e+6, "以字节为单位的清除阈值")
	flag.IntVar(&numItems, "count", 10000, "生成的文档数量")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.SetFlags(0)

	var (
		articles        []*Article
		countSuccessful uint64

		res *esapi.Response
		err error
	)
	addr := []string{`http://tencent.danny.games:9200/`}

	log.Printf(
		"\x1b[1mBulkIndexer\x1b[0m: documents [%s] workers [%d] flush [%s]",
		humanize.Comma(int64(numItems)), numWorkers, humanize.Bytes(uint64(flushBytes)))
	log.Println(strings.Repeat("▁", 65))

	// 使用第三方包来实现回退功能
	retryBackoff := backoff.NewExponentialBackOff()

	// 创建客户端。如果想使用最佳性能，请考虑使用第三方 http 包，benchmarks 示例中有写
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		// 429 太多请求
		RetryOnStatus: []int{502, 503, 504, 429},
		// 配置回退函数
		RetryBackoff: func(attempt int) time.Duration {
			if attempt == 1 {
				retryBackoff.Reset()
			}
			return retryBackoff.NextBackOff()
		},
		// 最多重试 5 次
		MaxRetries: 5,
		Addresses:  addr,
		Logger: &estransport.CurlLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	})
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 创建批量索引器。要使用最佳性能，考虑使用第三方 json 包，benchmarks 示例中有写
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         indexName,        // 默认索引名
		Client:        es,               // es 客户端
		NumWorkers:    numWorkers,       // 工作进程数量
		FlushBytes:    int(flushBytes),  // 清除上限
		FlushInterval: 30 * time.Second, // 定期清除间隔
	})
	if err != nil {
		log.Fatalf("Error creating the indexer: %s", err)
	}

	// 生成文章
	names := []string{"Alice", "John", "Mary"}
	for i := 1; i <= numItems; i++ {
		articles = append(articles, &Article{
			ID:        i,
			Title:     strings.Join([]string{"Title", strconv.Itoa(i)}, " "),
			Body:      "Lorem ipsum dolor sit amet...",
			Published: time.Now().Round(time.Second).UTC().AddDate(0, 0, i),
			Author: Author{
				FirstName: names[rand.Intn(len(names))],
				LastName:  "Smith",
			},
		})
	}
	log.Printf("→ Generated %s articles", humanize.Comma(int64(len(articles))))

	// 重新创建索引
	if res, err = es.Indices.Delete([]string{indexName}, es.Indices.Delete.WithIgnoreUnavailable(true)); err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}

	res.Body.Close()

	res, err = es.Indices.Create(indexName)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}
	res.Body.Close()

	start := time.Now().Local()

	// 循环收集
	for _, a := range articles {
		// 准备 data：序列化的 article
		data, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Cannot encode article %d: %s", a.ID, err)
		}

		// 添加 item 到 BulkIndexer
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				// Action 字段配置要执行的操作（索引、创建、删除、更新）
				Action: "index",
				// DocumentID 是文档 ID（可选）
				DocumentID: strconv.Itoa(a.ID),
				// Body 是 有效载荷的 io.Reader
				Body: bytes.NewReader(data),
				// OnSuccess 在每一个成功的操作后调用
				OnSuccess: func(c context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&countSuccessful, 1)
				},
				// OnFailure 在每一个失败的操作后调用
				OnFailure: func(c context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, e error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", biri.Error.Type, biri.Error.Reason)
					}
				},
			},
		)
		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}
	}

	// 关闭索引器
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	biStats := bi.Stats()

	// 报告结果：索引成功的文档的数量，错误的数量，耗时，索引速度
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if biStats.NumFailed > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			humanize.Comma(int64(biStats.NumFailed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	} else {
		log.Printf(
			"Sucessfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(biStats.NumFlushed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed))),
		)
	}
}
