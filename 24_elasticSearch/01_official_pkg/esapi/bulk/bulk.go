package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"math/rand"
	"os"

	"log"
	"strconv"
	"strings"
	"time"
)

// 批量索引 批量bulk

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
	_     = fmt.Print
	count int
	batch int
)

func init() {
	flag.IntVar(&count, "count", 1000, "生成的文档数量")
	flag.IntVar(&batch, "batch", 255, "每次发送的文档数量")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())
}

func main() {
	log.SetFlags(0)

	type bulkResponse struct {
		Errors bool `json:"errors"`
		Items  []struct {
			Index struct {
				ID     string `json:"_id"`
				Result string `json:"result"`
				Status int    `json:"status"`
				Error  struct {
					Type   string `json:"type"`
					Reason string `json:"reason"`
					Cause  struct {
						Type   string `json:"type"`
						Reason string `json:"reason"`
					} `json:"caused_by"`
				} `json:"error"`
			} `json:"index"`
		} `json:"items"`
	}

	var (
		buf bytes.Buffer
		res *esapi.Response
		err error
		raw map[string]interface{}
		blk *bulkResponse

		articles  []*Article
		indexName = "articles"

		numItems   int
		numErrors  int
		numIndexed int
		numBatches int
		currBatch  int
	)

	log.Printf(
		"\x1b[1mBulk\x1b[0m: documents [%s] batch size [%s]",
		humanize.Comma(int64(count)), humanize.Comma(int64(batch)))
	log.Println(strings.Repeat("_", 65))

	// 创建客户端
	addr := []string{`http://tencent.danny.games:9200/`}
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Logger: &estransport.CurlLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
		Addresses: addr,
	})
	if err != nil {
		panic(err)
	}

	// 生成文章
	names := []string{"Alice", "John", "Mary"}
	for i := 1; i < count+1; i++ {
		articles = append(articles, &Article{
			ID:        i,
			Title:     strings.Join([]string{"Title", strconv.Itoa(i)}, " "),
			Body:      "Lorem ipsum dolor sit amet...",
			Published: time.Now().Round(time.Second).Local().AddDate(0, 0, i),
			Author: Author{
				FirstName: names[rand.Intn(len(names))],
				LastName:  "Smith",
			},
		})
	}
	// 打印时人性化
	log.Printf("→ Generated %s articles", humanize.Comma(int64(len(articles))))
	fmt.Println("→ Sending batch ")

	// 重新创建索引
	if res, err = es.Indices.Delete([]string{indexName}); err != nil {
		log.Fatalf("Cannot delete index: %s", err)
	}

	res, err = es.Indices.Create(indexName)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}

	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}

	if count%batch == 0 {
		numBatches = count / batch
	} else {
		numBatches = count/batch + 1
	}

	start := time.Now().Local()

	// 循环收集
	for i, a := range articles {
		numItems++

		currBatch = i / batch
		if i == count-1 {
			currBatch++
		}

		// 准备元数据有效载荷
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, a.ID, "\n"))

		// 准备 data 有效载荷：序列化后的 article
		data, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Cannot encode article %d: %s", a.ID, err)
		}

		// 在 data 载荷中添加一个换行符
		data = append(data, "\n"...)

		// 将载荷添加到 buf 中
		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		// 达到阈值时，使用 buf 中的数据（请求体）执行 Bulk() 请求
		if i > 0 && i%batch == 0 || i == count-1 {
			fmt.Printf("[%d/%d] ", currBatch, numBatches)

			// 每 batch（本例中是255）个为一组发送
			res, err = es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
			if err != nil {
				log.Fatalf("Failur indexing batch %d: %s", currBatch, err)
			}

			// 如果整个请求失败，打印错误并标记所有文档都失败
			if res.IsError() {
				numErrors += numItems
				if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
					log.Fatalf("Failure to parse response body: %s", err)
				} else {
					log.Printf(" Error: [%d] %s: %s",
						res.StatusCode,
						raw["error"].(map[string]interface{})["type"],
						raw["error"].(map[string]interface{})["reason"],
					)
				}
			} else { // 一个成功的响应也可能因为一些特殊文档包含一些错误
				if err := json.NewDecoder(res.Body).Decode(&blk); err != nil {
					log.Fatalf("Failure to parse response body: %s", err)
				} else {
					for _, d := range blk.Items {
						// 对任何状态码大于 201 的请求进行处理
						if d.Index.Status > 201 {
							numErrors++
							log.Printf("  Error: [%d]: %s: %s: %s: %s",
								d.Index.Status,
								d.Index.Error.Type,
								d.Index.Error.Reason,
								d.Index.Error.Cause.Type,
								d.Index.Error.Cause.Reason,
							)
						} else {
							// 如果状态码小于等于 201，对成功的计数器 numIndexed 加 1
							numIndexed++
						}
					}
				}
			}

			// 关闭响应体，防止达到 Goroutines 或文件句柄限制
			res.Body.Close()

			// 重置 buf 和 items 计数器
			buf.Reset()
			numItems = 0
		}
	}

	// 报告结果：索引成功的文档的数量，错误的数量，耗时，索引速度
	fmt.Println()
	log.Println(strings.Repeat("▔", 65))

	dur := time.Since(start)

	if numErrors > 0 {
		log.Fatalf(
			"Indexed [%s] documents with [%s] errors in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			humanize.Comma(int64(numErrors)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	} else {
		log.Printf(
			"Successfuly indexed [%s] documents in %s (%s docs/sec)",
			humanize.Comma(int64(numIndexed)),
			dur.Truncate(time.Millisecond),
			humanize.Comma(int64(1000.0/float64(dur/time.Millisecond)*float64(numIndexed))),
		)
	}
}
