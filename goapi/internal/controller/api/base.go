package api

import (
	"clickhouse-api/internal/helper"
	"clickhouse-api/pkg/paginator"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

// swagger:response Response
type Response struct {
	// in: body
	Body struct {
		// Status of the error
		// in: bool
		Status bool `json:"status"`
		// Message of the error
		// in: string
		Message string `json:"message"`
	}
}

type Controller struct {
	Counter                *prometheus.GaugeVec
	ErrorCounter           *prometheus.GaugeVec
	Duration               *prometheus.HistogramVec
	RequestSize            *prometheus.HistogramVec
	LatencyPercentiles50th *prometheus.SummaryVec
	LatencyPercentiles90th *prometheus.SummaryVec
	LatencyPercentiles99th *prometheus.SummaryVec
}

func NewController(reg *prometheus.Registry) *Controller {
	counter := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "RequestCount",
		Name:      "request_count",
		Help:      "Количество запросов",
	}, []string{"url"})

	errorCounter := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "RequestErrorCount",
		Name:      "request_error_count",
		Help:      "Количество ошибок",
	}, []string{"status", "url"})

	duration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "RequestDuration",
		Name:      "request_duration",
		Help:      "Время выполнения запроса",
	}, []string{"url"})

	requestSize := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "RequestSize",
		Name:      "request_size",
		Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
		Help:      "Гистограмма размера запроса в байтах для отдельного URL",
	}, []string{"url"})

	latencyPercentiles50th := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "latencyPercentiles",
		Name:       "latency_percentiles_50th",
		Help:       "50th перцентиль времени выполнения запроса в секундах",
		Objectives: map[float64]float64{0.5: 0.05},
	}, []string{"url"})

	latencyPercentiles90th := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "latencyPercentiles",
		Name:       "latency_percentiles_90th",
		Help:       "90th перцентиль времени выполнения запроса в секундах",
		Objectives: map[float64]float64{0.9: 0.01},
	}, []string{"url"})

	latencyPercentiles99th := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "latencyPercentiles",
		Name:       "latency_percentiles_99th",
		Help:       "90th перцентиль времени выполнения запроса в секундах",
		Objectives: map[float64]float64{0.99: 0.001},
	}, []string{"url"})

	reg.MustRegister(
		counter,
		errorCounter,
		duration,
		requestSize,
		latencyPercentiles50th,
		latencyPercentiles90th,
		latencyPercentiles99th,
	)

	return &Controller{
		Counter:                counter,
		ErrorCounter:           errorCounter,
		Duration:               duration,
		RequestSize:            requestSize,
		LatencyPercentiles50th: latencyPercentiles50th,
		LatencyPercentiles90th: latencyPercentiles90th,
		LatencyPercentiles99th: latencyPercentiles99th,
	}
}

func (controller *Controller) Defer(request *http.Request, url string, startTime time.Time) {
	controller.UpdateMetrics(url, startTime)
	_ = request.Body.Close()
}

func (controller *Controller) UpdateMetrics(url string, startTime time.Time) {
	duration := time.Since(startTime).Seconds()
	controller.Duration.WithLabelValues(url).Observe(duration)
	controller.Counter.WithLabelValues(url).Dec()
	controller.LatencyPercentiles50th.WithLabelValues(url).Observe(duration)
	controller.LatencyPercentiles90th.WithLabelValues(url).Observe(duration)
	controller.LatencyPercentiles99th.WithLabelValues(url).Observe(duration)
}

func (controller *Controller) HandleError(url string, err error, writer http.ResponseWriter, httpStatus int) {
	controller.ErrorCounter.With(prometheus.Labels{"status": strconv.Itoa(httpStatus), "url": url}).Inc()
	log.Error(err)
	helper.Respond(writer, helper.Message(false, err.Error()), httpStatus)
}

func (controller *Controller) HandleValidationError(url string, message map[string]interface{}, writer http.ResponseWriter, httpStatus int) {
	controller.ErrorCounter.With(prometheus.Labels{"status": strconv.Itoa(httpStatus), "url": url}).Inc()
	log.Error(errors.New(fmt.Sprintf("%v", message)))
	helper.Respond(writer, message, httpStatus)
}

func (controller *Controller) PaginationResponse(writer http.ResponseWriter, paginator *paginator.Paginator, data interface{}) {
	response := helper.Message(true, "success")
	response["pagination"] = paginator
	response["data"] = data
	helper.Respond(writer, response, http.StatusOK)
}

func (controller *Controller) MessageResponse(writer http.ResponseWriter, message string, statusCode int) {
	helper.Respond(writer, helper.Message(true, message), statusCode)
}

func (controller *Controller) DataResponse(writer http.ResponseWriter, dataKey string, data interface{}, statusCode int) {
	response := map[string]interface{}{
		dataKey: data,
	}
	helper.Respond(writer, response, statusCode)
}
