package svc

import (
	"context"
	"fmt"
	"model_mall_backend/backend/internal/config"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// OrmHelper 封装了 GORM 的 PostgreSQL 数据库实例，主要用于连接管理和执行原生 SQL
type OrmHelper struct {
	db     *gorm.DB
	logger *LogHelper
	ctx    context.Context
}

// NewOrmHelper 创建一个新的 OrmHelper 实例
func NewOrmHelper(ctx context.Context, conf *config.Config, logger *LogHelper) (*OrmHelper, error) {
	// 构建 PostgreSQL DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		conf.PostgreSQL.Host,
		conf.PostgreSQL.Username,
		conf.PostgreSQL.Password,
		conf.PostgreSQL.Database,
		conf.PostgreSQL.Port,
		conf.PostgreSQL.SSLMode,
	)

	// 使用 GORM 配置连接
	gormConfig := &gorm.Config{
		// SkipDefaultTransaction: true, // <--- 注释掉此项
		// 解释：GORM 默认为写入操作启用事务以保证一致性。
		// 禁用它（设为 true）主要是为了微小的性能优化或特殊需求。
		// 对于大多数情况，保留默认（false）是推荐且安全的做法。
		// 因此，此行被注释，意味着使用 GORM 的默认事务行为。

		// 禁用自动时间追踪（如果不需要 CreateAt/UpdateAt）
		// NowFunc: func() time.Time { return time.Now() }, // 可以自定义时间函数
		// Logger: ... // 可配置为自定义 Logger 或使用默认
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		logger.Error(err, "Failed to connect to postgres")
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	// 获取底层的 database/sql.DB 对象以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		logger.Error(err, "Failed to get database/sql.DB from GORM")
		return nil, fmt.Errorf("failed to get database/sql.DB from GORM: %v", err)
	}

	// 配置连接池
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	// 测试连接
	sqlCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(sqlCtx); err != nil {
		logger.Error(err, "Failed to ping postgres")
		return nil, fmt.Errorf("failed to ping postgres: %v", err)
	}

	return &OrmHelper{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}, nil
}

// Close 关闭数据库连接池
func (p *OrmHelper) Close() {
	if p.db != nil {
		sqlDB, err := p.db.DB()
		if err != nil {
			p.logger.Error(err, "Failed to get *sql.DB from GORM for closing")
			return
		}
		sqlDB.Close()
	}
}

// GetDB 获取 GORM DB 实例 (可用于 ORM 操作或获取 *sql.DB)
func (p *OrmHelper) GetDB() *gorm.DB {
	return p.db
}

// ExecContext 执行非查询 SQL 语句 (如 INSERT, UPDATE, DELETE)
func (p *OrmHelper) ExecContext(ctx context.Context, sql string, args ...interface{}) (int64, error) {
	start := time.Now()

	result := p.db.WithContext(ctx).Exec(sql, args...)
	err := result.Error

	duration := time.Since(start)
	if duration > time.Second {
		p.logger.Slow("Slow SQL execution (took %v): %s", duration, sql)
	}

	if err != nil {
		p.logger.Error(err, "Failed to execute SQL: %s", sql)
		return 0, err
	}

	rowsAffected := result.RowsAffected
	p.logger.Debug("SQL executed successfully: %s, RowsAffected: %d", sql, rowsAffected)
	return rowsAffected, nil
}

// QueryRowContext 查询单行数据
// dest 应该是指向结构体或基本类型的指针，GORM 会将结果扫描到其中
func (p *OrmHelper) QueryRowContext(ctx context.Context, dest interface{}, sql string, args ...interface{}) error {
	start := time.Now()

	// 使用 Raw + Scan 执行原生 SQL 查询并将结果扫描到 dest
	err := p.db.WithContext(ctx).Raw(sql, args...).Scan(dest).Error

	duration := time.Since(start)
	if duration > time.Second {
		p.logger.Slow("Slow SQL query row (took %v): %s", duration, sql)
	}

	if err != nil {
		// 区分 "未找到记录" 错误和其它错误
		if err == gorm.ErrRecordNotFound {
			p.logger.Debug("No record found for SQL query: %s", sql)
		} else {
			p.logger.Error(err, "Failed to execute SQL query row: %s", sql)
		}
		return err
	}

	p.logger.Debug("SQL query row executed successfully: %s", sql)
	return nil
}

// QueryContext 查询多行数据
// dest 应该是指向结构体切片的指针，GORM 会将结果扫描到其中
func (p *OrmHelper) QueryContext(ctx context.Context, dest interface{}, sql string, args ...interface{}) error {
	start := time.Now()

	// 使用 Raw + Scan 执行原生 SQL 查询并将结果扫描到 dest (切片)
	err := p.db.WithContext(ctx).Raw(sql, args...).Scan(dest).Error

	duration := time.Since(start)
	if duration > time.Second {
		p.logger.Slow("Slow SQL query rows (took %v): %s", duration, sql)
	}

	if err != nil {
		p.logger.Error(err, "Failed to execute SQL query rows: %s", sql)
		return err
	}

	p.logger.Debug("SQL query rows executed successfully: %s", sql)
	return nil
}

// Transaction executes a function within a database transaction.
func (p *OrmHelper) Transaction(fc func(tx *OrmHelper) error) error {
	return p.db.Transaction(func(tx *gorm.DB) error {
		txHelper := &OrmHelper{
			db:     tx,
			logger: p.logger,
			ctx:    p.ctx,
		}
		return fc(txHelper)
	})
}
