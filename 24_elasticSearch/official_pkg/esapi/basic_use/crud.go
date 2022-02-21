// $ go run _examples/main.go

package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

// 服务器使用的是v7版本

func initClient() (err error) {
	// 配置服务器地址
	addr := []string{`http://tencent.danny.games:9200/`}
	// 配置http数据传输
	transport := &http.Transport{
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second,
	}
	config := elasticsearch.Config{Addresses: addr, Transport: transport}
	es, err = elasticsearch.NewClient(config)
	return

}

var es *elasticsearch.Client

func getClusterInfo() {
	var (
		r map[string]interface{}
	)

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	// 关闭是为了复用tcp🔗
	//  It is critical to both close the response body and to consume it, in order to re-use persistent TCP connections in the default HTTP transport
	defer res.Body.Close()
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print version number.
	// Print client and server version numbers.
	log.Printf("Client: %s", elasticsearch.Version)                           // Client: 7.16.0
	log.Printf("Server: %s", r["version"].(map[string]interface{})["number"]) // Server: 7.13.2
	log.Println(strings.Repeat("~", 37))
	/*
		json返回内容
			{
			  "name" : "danny",
			  "cluster_name" : "elasticsearch",
			  "cluster_uuid" : "pr0k5JSSTXq2v4CB8OLdQA",
			  "version" : {
			    "number" : "7.13.2",
			    "build_flavor" : "default",
			    "build_type" : "tar",
			    "build_hash" : "4d960a0733be83dd2543ca018aa4ddc42e956800",
			    "build_date" : "2021-06-10T21:01:55.251515791Z",
			    "build_snapshot" : false,
			    "lucene_version" : "8.8.2",
			    "minimum_wire_compatibility_version" : "6.8.0",
			    "minimum_index_compatibility_version" : "6.0.0-beta1"
			  },
			  "tagline" : "You Know, for Search"
			}
	*/
}

func insertIndex() {
	var wg sync.WaitGroup
	for i, title := range []string{"Test One", "Test Two"} {
		wg.Add(1)

		go func(i int, title string) {
			defer wg.Done()

			// Set up the request object directly.
			req := esapi.IndexRequest{
				Index:      "test",
				DocumentID: strconv.Itoa(i + 1),
				Body:       strings.NewReader(`{"title" : "` + title + `"}`),
				Refresh:    "true", // 我们设置 Refresh 为 true。这在实际的使用中并不建议，原因是每次写入的时候都会 refresh。当我们面对大量的数据时，这样的操作会造成效率的底下。
			}

			// Perform the request with the client.
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}
		}(i, title)
	}
	wg.Wait()

	log.Println(strings.Repeat("-", 37))
}

func searchIndex() {
	var (
		r map[string]interface{}
	)
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithBody(strings.NewReader(`{"query" : { "match" : { "title" : "test" } }}`)),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	}

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// Print the response status, number of results, and request duration.
	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)
	// Print the ID and document source for each hit.
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		log.Printf(" * ID=%s, %s", hit.(map[string]interface{})["_id"], hit.(map[string]interface{})["_source"])
	}

	log.Println(strings.Repeat("=", 37))
}
