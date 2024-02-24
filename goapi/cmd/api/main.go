//	 Clickhouse API:
//	  version: 1.0
//	  title: Clickhouse API
//	 Schemes: http
//	 Host:
//	 BasePath: /api/v1
//		Consumes:
//		- application/json
//	 Produces:
//	 - application/json
//	 SecurityDefinitions:
//	  basicAuth:
//	   type: basic
//	 swagger:meta
package main

import (
	"clickhouse-api/internal/repository/base/clickhouse"
	"clickhouse-api/internal/repository/test_data"
	"clickhouse-api/internal/service/config"
	"clickhouse-api/internal/service/logger"
	clickhouseQueue "clickhouse-api/internal/service/queue"
	"clickhouse-api/internal/service/router"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
)

func main() {
	mainConfig := config.New()
	reg := prometheus.NewRegistry()
	connectionMetrics := clickhouse.NewConnectionMetrics(reg)
	testDataRepository := test_data.New(mainConfig, connectionMetrics)
	createQueue := clickhouseQueue.New(mainConfig, "create", test_data.Table, testDataRepository)
	updateQueue := clickhouseQueue.New(mainConfig, "update", test_data.Table, testDataRepository)
	deleteQueue := clickhouseQueue.New(mainConfig, "delete", test_data.Table, testDataRepository)
	defer testDataRepository.Close()
	metricsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	appRouter := router.New(mainConfig, createQueue, updateQueue, deleteQueue, testDataRepository, reg, metricsHandler)
	go ListenQueue(createQueue)
	go ListenQueue(updateQueue)
	go ListenQueue(deleteQueue)
	logger.Init(mainConfig)

	err := http.ListenAndServe(mainConfig.Application.Host+":"+strconv.Itoa(mainConfig.Application.Port), appRouter.GetRouter())
	if err != nil {
		log.Fatalln(fmt.Errorf("ListenAndServe: %+v", err))
	}
}

func ListenQueue(queue *clickhouseQueue.FileQueue) {
	queue.Listen()
}
