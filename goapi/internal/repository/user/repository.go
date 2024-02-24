package user

import (
	"clickhouse-api/internal/repository/base/sqlite"
	"clickhouse-api/internal/service/config"
	"fmt"
)

const (
	repositoryName = "userRepository "
)

type Repository struct {
	*sqlite.Repository
}

func New(conf *config.Config) *Repository {
	return &Repository{
		sqlite.New(conf.Sqlite),
	}
}

func (repository *Repository) GetCountByLoginAndPassword(login, token string) (uint64, error) {
	const countError = repositoryName + "get count by login and password error: "
	var count uint64 = 0
	if err := repository.Connection.Get(&count, "SELECT COUNT(*) FROM main.users WHERE login=$1 AND token=$2", login, token); err != nil {
		return 0, fmt.Errorf(countError+"get query error: %+v", err)
	}
	return count, nil
}
