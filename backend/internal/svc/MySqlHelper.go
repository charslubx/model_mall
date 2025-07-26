package svc

import (
	"database/sql"
	"fmt"
	"model_mall_backend/backend/internal/config"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type MySqlHelper struct {
	conn   sqlx.SqlConn
	logger *LogHelper
}

func NewMySqlHelper(conf *config.Config, logger *LogHelper) (*MySqlHelper, error) {
	// go-zero 的数据源格式：user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai
	dataSource := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Asia%%2FShanghai",
		conf.MySQL.Username,
		conf.MySQL.Password,
		conf.MySQL.Host,
		conf.MySQL.Port,
		conf.MySQL.Database,
	)

	// 创建连接
	conn := sqlx.NewMysql(dataSource)

	return &MySqlHelper{
		conn:   conn,
		logger: logger,
	}, nil
}

// Exec 执行写操作
func (m *MySqlHelper) Exec(query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := m.conn.Exec(query, args...)
	duration := time.Since(start)

	if duration > time.Second {
		m.logger.Slow("Slow MySQL execution (took %v): %s", duration, query)
	}

	if err != nil {
		m.logger.Error(err, "Failed to execute MySQL query: %s", query)
		return nil, err
	}

	m.logger.Debug("MySQL executed successfully: %s", query)
	return result, nil
}

// QueryRow 查询单行
func (m *MySqlHelper) QueryRow(v interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := m.conn.QueryRow(v, query, args...)
	duration := time.Since(start)

	if duration > time.Second {
		m.logger.Slow("Slow MySQL query (took %v): %s", duration, query)
	}

	if err != nil {
		m.logger.Error(err, "Failed to query MySQL row: %s", query)
		return err
	}

	m.logger.Debug("MySQL query row successful: %s", query)
	return nil
}

// QueryRows 查询多行
func (m *MySqlHelper) QueryRows(v interface{}, query string, args ...interface{}) error {
	start := time.Now()
	err := m.conn.QueryRows(v, query, args...)
	duration := time.Since(start)

	if duration > time.Second {
		m.logger.Slow("Slow MySQL query (took %v): %s", duration, query)
	}

	if err != nil {
		m.logger.Error(err, "Failed to query MySQL rows: %s", query)
		return err
	}

	m.logger.Debug("MySQL query rows successful: %s", query)
	return nil
}

// Transact 事务操作
func (m *MySqlHelper) Transact(fn func(session sqlx.Session) error) error {
	return m.conn.Transact(fn)
}

// GetConn 获取原始连接
func (m *MySqlHelper) GetConn() sqlx.SqlConn {
	return m.conn
}
