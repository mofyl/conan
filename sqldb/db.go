package sqldb

import (
	"context"
	"sync"
)

type ExecResult interface {
	LastInsertId() (int64, error)

	RowsAffected() (int64, error)
}

type SqlDB interface {
	//getConnStr() string
	Exec(ctx context.Context, sql string, args ...interface{}) (ExecResult, error)
	QueryOne(ctx context.Context, res interface{}, sql string, args ...interface{}) error
	QueryAll(ctx context.Context, res interface{}, sql string, args ...interface{}) error
	Count(ctx context.Context, res interface{}, sql string, args ...interface{}) error
}

type baseConf struct {
}

type DriverType string

var (
	MYSQL DriverType

	mysqlDriver *mysqlDB
	mysqlOnce   *sync.Once
)

func init() {
	mysqlOnce = &sync.Once{}
}

func GetDB(driverType DriverType) SqlDB {
	switch driverType {
	default:
		mysqlOnce.Do(func() {
			mysqlDriver = newMysqlDriver()
		})
		return mysqlDriver
	}
}
