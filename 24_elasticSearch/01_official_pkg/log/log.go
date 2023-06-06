package main

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/estransport"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(0)

	var es *elasticsearch.Client

	addr := []string{`http://tencent.danny.games:9200/`}
	//es, _ = elasticsearch.NewClient(elasticsearch.Config{
	//	Logger:    &estransport.TextLogger{Output: os.Stdout},
	//	Addresses: addr,
	//})
	//run(es, "Text")
	//
	//es, _ = elasticsearch.NewClient(elasticsearch.Config{
	//	Logger: &estransport.ColorLogger{Output: os.Stdout},
	//})
	//run(es, "Color")
	//
	//es, _ = elasticsearch.NewClient(elasticsearch.Config{
	//	Logger: &estransport.ColorLogger{
	//		Output:             os.Stdout,
	//		EnableRequestBody:  true,
	//		EnableResponseBody: true,
	//	},
	//	Addresses: addr,
	//})
	//run(es, "Request/Response body")

	es, _ = elasticsearch.NewClient(elasticsearch.Config{
		Logger: &estransport.CurlLogger{
			Output:             os.Stdout,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
		Addresses: addr,
	})
	run(es, "Curl")

	//es, _ = elasticsearch.NewClient(elasticsearch.Config{
	//	Logger: &estransport.JSONLogger{
	//		Output: os.Stdout,
	//	},
	//	Addresses: addr,
	//})
	//run(es, "JSON")
}

func run(es *elasticsearch.Client, name string) {
	log.Println("███", fmt.Sprintf("\x1b[1m%s\x1b[0m", name), strings.Repeat("█", 75-len(name)))

	es.Delete("test", "1")
	es.Exists("test", "1")

	es.Index(
		"test",
		strings.NewReader(`{"title": "logging"}`),
		es.Index.WithRefresh("true"),
		es.Index.WithPretty(),
		es.Index.WithFilterPath("result", "_id"),
	)

	es.Search(es.Search.WithQuery("{FAIL"))

	res, err := es.Search(
		es.Search.WithIndex("test"),
		es.Search.WithBody(strings.NewReader(`{"query": {"match": {"title": "logging"}}}`)),
		es.Search.WithSize(1),
		es.Search.WithPretty(),
		es.Search.WithFilterPath("took", "hits.hits"),
	)

	s := res.String()

	if len(s) <= len("[200 OK] ") {
		log.Fatal("Response body is empty")
	}

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	log.Println()
}
