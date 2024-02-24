package repository

type Interface interface {
	InitConnection() error
	CloseConnection()
}
