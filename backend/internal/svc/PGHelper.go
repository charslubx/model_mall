package svc

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"model_mall_backend/backend/internal/config"
	"time"
)

type PGHelper struct {
	pool   *pgxpool.Pool
	logger *LogHelper
	ctx    context.Context
}

func NewPGHelper(ctx context.Context, conf *config.Config, logger *LogHelper) (*PGHelper, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		conf.PostgreSQL.Host,
		conf.PostgreSQL.Port,
		conf.PostgreSQL.Username,
		conf.PostgreSQL.Password,
		conf.PostgreSQL.Database,
		conf.PostgreSQL.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error(err, "Failed to parse postgres config")
		return nil, fmt.Errorf("failed to parse postgres config: %v", err)
	}

	// 设置连接池配置
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		logger.Error(err, "Failed to connect to postgres")
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	return &PGHelper{
		pool:   pool,
		logger: logger,
		ctx:    ctx,
	}, nil
}

// Close 关闭数据库连接池
func (p *PGHelper) Close() {
	if p.pool != nil {
		p.pool.Close()
	}
}

// GetPool 获取连接池实例
func (p *PGHelper) GetPool() *pgxpool.Pool {
	return p.pool
}

// ExecContext 执行SQL语句
func (p *PGHelper) ExecContext(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	start := time.Now()
	result, err := p.pool.Exec(ctx, sql, args...)
	duration := time.Since(start)

	if duration > time.Second {
		p.logger.Slow("Slow SQL execution (took %v): %s", duration, sql)
	}

	if err != nil {
		p.logger.Error(err, "Failed to execute SQL: %s", sql)
		return pgconn.CommandTag{}, err
	}

	p.logger.Debug("SQL executed successfully: %s", sql)
	return result, nil
}

// QueryRowContext 查询单行数据
func (p *PGHelper) QueryRowContext(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	start := time.Now()
	row := p.pool.QueryRow(ctx, sql, args...)
	duration := time.Since(start)

	if duration > time.Second {
		p.logger.Slow("Slow SQL query (took %v): %s", duration, sql)
	}

	p.logger.Debug("SQL query executed: %s", sql)
	return row
}

// QueryContext 查询多行数据
func (p *PGHelper) QueryContext(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	start := time.Now()
	rows, err := p.pool.Query(ctx, sql, args...)
	duration := time.Since(start)

	if duration > time.Second {
		p.logger.Slow("Slow SQL query (took %v): %s", duration, sql)
	}

	if err != nil {
		p.logger.Error(err, "Failed to execute SQL query: %s", sql)
		return nil, err
	}

	p.logger.Debug("SQL query executed successfully: %s", sql)
	return rows, nil
}
