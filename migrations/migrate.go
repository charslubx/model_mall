package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DatabaseConfig 数据库配置结构
type DatabaseConfig struct {
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Username string `yaml:"Username"`
	Password string `yaml:"Password"`
	Database string `yaml:"Database"`
	SSLMode  string `yaml:"SSLMode"`
}

// Config 配置结构
type Config struct {
	PostgreSQL DatabaseConfig `yaml:"PostgreSQL"`
}

func main() {
	// 加载配置
	configPath := "../backend/etc/backend-api.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	c, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接数据库
	db, err := connectDB(c)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	fmt.Println("开始执行数据库迁移...")

	// 执行SQL文件迁移
	if err := runSQLMigrations(db); err != nil {
		log.Fatalf("SQL迁移失败: %v", err)
	}

	fmt.Println("数据库迁移完成！")
}

// 加载配置文件
func loadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
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
