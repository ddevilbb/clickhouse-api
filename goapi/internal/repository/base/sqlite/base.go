package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

const (
	DriverName = "sqlite"
)

type Config struct {
	ConnectionFileSrc string
}

type Repository struct {
	Connection *sqlx.DB
}

func New(config *Config) *Repository {
	connection, err := sqlx.Open("sqlite3", config.ConnectionFileSrc)
	if err != nil {
		log.Errorf("%+v", err)
		return nil
	}

	return &Repository{
		Connection: connection,
	}
}
