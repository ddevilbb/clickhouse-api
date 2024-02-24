package logger

import (
	"clickhouse-api/internal/service/config"
	"clickhouse-api/pkg/elk_writer"
	log "github.com/sirupsen/logrus"
)

type jsonLogger struct {
	Type        string
	Environment string
	formatter   log.Formatter
}

func Init(conf *config.Config) {
	log.SetFormatter(jsonLogger{
		Type:        "clickhouse-api",
		Environment: conf.Application.Environment,
		formatter:   &log.JSONFormatter{},
	})
	log.SetOutput(elk_writer.New(conf.ElkWriter))
	log.SetReportCaller(true)
}

func (l jsonLogger) Format(entry *log.Entry) ([]byte, error) {
	entry.Data["environment"] = l.Environment
	entry.Data["type"] = l.Type

	return l.formatter.Format(entry)
}
