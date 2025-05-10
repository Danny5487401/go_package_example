package main

import (
	"context"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
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
	switch results.(type) {
	case model.Vector:
		for _, result := range results.(model.Vector) {
			log.Printf("result:\n timestamp:%v  metric:%v \n values:%v \n", result.Timestamp.Time(), result.Metric, result.Value)
		}
		/*
			2025/03/30 11:16:08 result:
			 timestamp:2025-03-30 11:16:08.279 +0800 CST  metric:instance:node_cpu_utilisation:rate5m{container="node-exporter", endpoint="http-metrics", instance="172.16.7.33:9100", job="node-exporter", namespace="monitor", pod="prometheus-prometheus-node-exporter-frmx4", service="prometheus-prometheus-node-exporter"}
			 values:0.07596296296295435
		*/
	default:
		log.Printf("error: The Prometheus results should not be type: %v.\n", results.Type())
	}

}
