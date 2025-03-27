package main

import (
	"context"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"log"
	"time"
)

func main() {
	promAddress := "http://prometheus-kube-prometheus-prometheus.monitor.svc.cluster.local:9090"
	client, err := api.NewClient(api.Config{
		Address: promAddress,
	})
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	promQuery := "instance:node_cpu_utilisation:rate5m"
	results, warnings, err := v1api.Query(ctx, promQuery, time.Now())
	if err != nil {
		log.Printf("err: %v\n", err)
		return
	}
	if len(warnings) > 0 {
		log.Printf("Warnings: %v\n", warnings)
	}
	log.Printf("result:\n type:%v \n values:%v \n", results.Type(), results.String())

}
