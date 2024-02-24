package test_data

import (
	"clickhouse-api/internal/controller/api"
	"clickhouse-api/internal/repository/test_data"
)

type Readable struct {
	Repository *test_data.Repository
	*api.Controller
}
