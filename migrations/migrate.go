package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config 迁移工具最小配置（从 ../backend/etc/backend-api.yaml 加载）
type Config struct {
	PostgreSQL struct {
		Host     string
		Port     int
		Username string
		Password string
		Database string
		SSLMode  string
	}
}

func main() {
	// 加载配置
	var c Config
	conf.MustLoad("../backend/etc/backend-api.yaml", &c)

	// 连接数据库
	db, err := connectDB(&c)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	fmt.Println("开始执行数据库迁移...")

	// 方式1：执行SQL文件迁移
	if err := runSQLMigrations(db); err != nil {
		log.Fatalf("SQL迁移失败: %v", err)
	}

	fmt.Println("数据库迁移完成！")
}

// 连接数据库
func connectDB(c *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.PostgreSQL.Host,
		c.PostgreSQL.Username,
		c.PostgreSQL.Password,
		c.PostgreSQL.Database,
		c.PostgreSQL.Port,
		c.PostgreSQL.SSLMode,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
}

// 执行SQL文件迁移
func runSQLMigrations(db *gorm.DB) error {
	// 获取所有SQL文件
	files, err := filepath.Glob("*.sql")
	if err != nil {
		return fmt.Errorf("获取SQL文件失败: %v", err)
	}

	// 按文件名排序
	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("执行迁移文件: %s\n", file)

		// 读取SQL文件内容
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("读取文件 %s 失败: %v", file, err)
		}

		// 分割SQL语句（以分号分隔）
		sqlStatements := strings.Split(string(content), ";")

		for _, stmt := range sqlStatements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" || strings.HasPrefix(stmt, "--") {
				continue
			}

			// 执行SQL语句
			if err := db.Exec(stmt).Error; err != nil {
				// 忽略已存在的表/索引等错误
				if strings.Contains(err.Error(), "already exists") {
					fmt.Printf("跳过已存在的对象: %s\n", stmt[:min(50, len(stmt))])
					continue
				}
				return fmt.Errorf("执行SQL失败 [%s]: %v", file, err)
			}
		}

		fmt.Printf("✓ %s 执行完成\n", file)
	}

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
