package sqldb

import (
	"context"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
	"conan/config"
)

type mysqlConf struct {
	Addr         string `yaml:"addr"`
	UserName     string `yaml:"userName"`
	Pwd          string `yaml:"pwd"`
	DBName       string `yaml:"dbName"`
	MaxOpenConn  int    `yaml:"maxOpenConn"`
	MaxIdleConn  int    `yaml:"maxIdleConn"`
	ConnLiftTime int    `yaml:"connLiftTime"`
}

const (
	EXPTIME = 3 // 单位是s
)

type mysqlDB struct {
	engine *sqlx.DB
	conf   *mysqlConf
}

func newMysqlDriver() *mysqlDB {

	db := &mysqlDB{}
	conf := &mysqlConf{}
	err := config.Decode(config.CONFIG_FILE_NAME, "mysql", conf)
	if err != nil {
		panic("Create Mysql Driver Fail " + err.Error())
	}
	db.conf = conf
	var engine *sqlx.DB
	engine, err = db.Connect()
	if err != nil {
		panic("Mysql Connect Fail " + err.Error())
	}

	engine.SetConnMaxLifetime(time.Duration(db.conf.ConnLiftTime) * time.Second)
	engine.SetMaxIdleConns(db.conf.MaxIdleConn)
	engine.SetMaxOpenConns(db.conf.MaxOpenConn)

	db.engine = engine

	return db
}

func (m *mysqlDB) Connect() (*sqlx.DB, error) {
	connStr := m.getConnStr()
	return sqlx.Connect("mysql", connStr)

}

func (m *mysqlDB) getConnStr() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", m.conf.UserName, m.conf.Pwd, m.conf.Addr, m.conf.DBName)
}

func (m *mysqlDB) QueryOne(ctx context.Context, res interface{}, sql string, args ...interface{}) error {
	row, err := m.engine.QueryxContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	if row.Next() {
		return row.StructScan(res)
	}

	return nil

}

func (m *mysqlDB) QueryAll(ctx context.Context, res interface{}, sql string, args ...interface{}) error {
	return m.engine.SelectContext(ctx, res, sql, args...)
}

func (m *mysqlDB) Exec(ctx context.Context, sql string, args ...interface{}) (ExecResult, error) {
	return m.engine.ExecContext(ctx, sql, args...)
}

func (m *mysqlDB) Count(ctx context.Context, res interface{}, sql string, args ...interface{}) error {
	return m.engine.GetContext(ctx, res, sql, args...)
}
