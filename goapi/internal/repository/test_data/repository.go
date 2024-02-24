package test_data

import (
	"clickhouse-api/internal/helper"
	"clickhouse-api/internal/model"
	"clickhouse-api/internal/repository/base/clickhouse"
	"clickhouse-api/internal/service/config"
	"fmt"
	"strconv"
)

const (
	Table          = "test_data"
	repositoryName = "testDataRepository "
)

type Repository struct {
	*clickhouse.Repository
}

func New(conf *config.Config, metrics *clickhouse.ConnectionMetrics) *Repository {
	return &Repository{
		clickhouse.New(conf.Clickhouse, metrics),
	}
}

func (repository *Repository) BulkInsert(list []string) error {
	const bulkInsertError = repositoryName + "bulk insert error: "
	transaction, err := repository.Connection.Begin()
	if err != nil {
		return fmt.Errorf(bulkInsertError+"transaction begin error: %+v", err)
	}
	sql := fmt.Sprintf(
		"INSERT INTO %s (sign, version, `data`, created_at) VALUES (?, ?, ?, ?) "+
			"SETTINGS async_insert = 1, async_insert_busy_timeout_ms = 5000, async_insert_stale_timeout_ms = 5000",
		Table,
	)
	stmt, err := transaction.Prepare(sql)
	if err != nil {
		return fmt.Errorf(bulkInsertError+"transaction prepare query error: %+v", err)
	}
	for _, item := range list {
		if item == "" {
			continue
		}
		testData, err := model.CreateTestDataFromJsonString(item)
		if err != nil {
			return fmt.Errorf(bulkInsertError+"%+v", err)
		}
		_, err = stmt.Exec(
			testData.Sign,
			testData.Version,
			testData.Data,
			testData.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf(bulkInsertError+"transaction exec error: %+v", err)
		}
	}
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf(bulkInsertError+"transaction commit error: %+v", err)
	}
	return nil
}

func (repository *Repository) BulkUpdate(list []string) error {
	return repository.bulkInsertWithId(list)
}

func (repository *Repository) BulkDelete(list []string) error {
	return repository.bulkInsertWithId(list)
}

func (repository *Repository) bulkInsertWithId(list []string) error {
	transaction, err := repository.Connection.Begin()
	if err != nil {
		return err
	}
	sql := fmt.Sprintf(
		"INSERT INTO %s (id, sign, version, `data`, created_at) VALUES (?, ?, ?, ?, ?) "+
			"SETTINGS async_insert = 1, async_insert_busy_timeout_ms = 5000, async_insert_stale_timeout_ms = 5000",
		Table,
	)
	stmt, err := transaction.Prepare(sql)
	if err != nil {
		return err
	}
	for _, item := range list {
		if item == "" {
			continue
		}
		testData, err := model.CreateTestDataFromJsonString(item)
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(
			testData.Id,
			testData.Sign,
			testData.Version,
			testData.Data,
			testData.CreatedAt,
		); err != nil {
			return err
		}
	}
	err = transaction.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (repository *Repository) Count() (uint64, error) {
	var count uint64 = 0
	const countError = repositoryName + "count error: "
	sql := fmt.Sprintf(`SELECT COUNT(*) FROM %s FINAL`, Table)
	if err := repository.Connection.Get(&count, sql); err != nil {
		return 0, fmt.Errorf(countError+"get query error: %+v", err)
	}

	return count, nil
}

func (repository *Repository) CountBy(params map[string][]string) (uint64, error) {
	var count uint64 = 0
	whereSql, err := helper.PrepareSqlWhere(params)
	if err != nil {
		return 0, err
	}
	sql := fmt.Sprintf(
		`SELECT SUM(sign) FROM %s %s`,
		Table,
		whereSql,
	)
	if err = repository.Connection.Get(&count, sql); err != nil {
		return 0, err
	}

	return count, nil
}

func (repository *Repository) FindIdsBy(offset, limit uint64, params map[string][]string) ([]string, error) {
	var result []string
	whereSql, err := helper.PrepareSqlWhere(params)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf(
		`SELECT "id" FROM %s %s GROUP BY id, version HAVING SUM(sign) > 0 LIMIT %s, %s`,
		Table,
		whereSql,
		strconv.FormatUint(offset, 10),
		strconv.FormatUint(limit, 10),
	)
	if err := repository.Connection.Select(&result, sql); err != nil {
		return nil, err
	}

	return result, nil
}

func (repository *Repository) FindByIds(ids []string, orderParams map[string]string) ([]model.TestData, error) {
	var items []model.TestData
	params := map[string][]string{"id": ids}
	whereSql, err := helper.PrepareSqlWhere(params)
	if err != nil {
		return nil, err
	}
	orderBy, err := helper.PrepareSqlOrder(orderParams)
	if err != nil {
		return nil, err
	}
	sql := fmt.Sprintf(
		`SELECT id, sign, version, "data", created_at FROM %s %s %s`,
		Table,
		whereSql,
		orderBy,
	)
	if err := repository.Connection.Select(&items, sql); err != nil {
		return nil, err
	}

	return items, nil
}

func (repository *Repository) Find(id string) (model.TestData, error) {
	var event model.TestData

	sql := fmt.Sprintf(
		`SELECT id, sign, version, "data", created_at FROM %s FINAL WHERE "id"= ?`,
		Table,
	)

	if err := repository.Connection.Get(&event, sql, id); err != nil {
		return model.TestData{}, err
	}

	return event, nil
}
