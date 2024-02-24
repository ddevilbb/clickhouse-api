package clickhouse

import (
	"github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

const (
	DriverName = "clickhouse"
)

type Config struct {
	ConnectionUrl string
}

type Repository struct {
	Connection *sqlx.DB
	Metrics    *ConnectionMetrics
}

type ConnectionMetrics struct {
	activeConnections prometheus.Gauge
}

func NewConnectionMetrics(reg *prometheus.Registry) *ConnectionMetrics {
	activeConnections := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "ClickhouseApi_Clickhouse",
		Name:      "active_connections",
		Help:      "количество активных соединений с кликхаусом",
	})

	reg.MustRegister(activeConnections)

	return &ConnectionMetrics{
		activeConnections: activeConnections,
	}
}

func New(config *Config, metrics *ConnectionMetrics) *Repository {
	connection, err := sqlx.Open(DriverName, config.ConnectionUrl)
	if err != nil {
		log.Errorf("%+v", err)
		return nil
	}
	if err = connection.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			log.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			return nil
		} else {
			log.Errorf("clickhouse connection ping error: %+v", err)
			return nil
		}
	}

	metrics.activeConnections.Inc()

	return &Repository{
		Connection: connection,
		Metrics:    metrics,
	}
}

func (repo *Repository) Close() {
	if repo.Connection != nil {
		repo.Close()
		repo.Metrics.activeConnections.Dec()
	}
}
