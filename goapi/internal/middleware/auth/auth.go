package auth

import (
	"clickhouse-api/internal/repository/user"
	"clickhouse-api/internal/service/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
	"strings"
)

type Middleware struct {
	repository         *user.Repository
	authSuccessCounter prometheus.Gauge
	authFailureCounter prometheus.Gauge
}

func New(conf *config.Config, reg *prometheus.Registry) *Middleware {
	authSuccessCounter := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "ClickhouseApi_Authentication",
		Name:      "success_auth_total",
		Help:      "Количество успешных попыток аутентификации",
	})

	authFailureCounter := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "ClickhouseApi_Authentication",
		Name:      "failure_auth_total",
		Help:      "Количество неудачных попыток аутентификации",
	})

	reg.MustRegister(authSuccessCounter, authFailureCounter)

	return &Middleware{
		repository:         user.New(conf),
		authSuccessCounter: authSuccessCounter,
		authFailureCounter: authFailureCounter,
	}
}

func (auth *Middleware) Basic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		login, token, ok := request.BasicAuth()
		if !ok || len(strings.TrimSpace(login)) < 1 || len(strings.TrimSpace(token)) < 1 {
			writer.WriteHeader(http.StatusUnauthorized)
			auth.authFailureCounter.Inc()
			return
		}
		if count, err := auth.repository.GetCountByLoginAndPassword(login, token); count == 0 || err != nil {
			writer.WriteHeader(http.StatusUnauthorized)
			auth.authFailureCounter.Inc()
			return
		}

		auth.authSuccessCounter.Inc()
		next.ServeHTTP(writer, request)
	})
}
