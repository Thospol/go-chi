package sql

import (
	"fmt"
	"time"

	"saaa-api/internal/core/config"

	"gorm.io/driver/mysql"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	// Database global variable database
	Database = &gorm.DB{}
)

// MysqlConfig config posgresql
type MysqlConfig struct {
	Database *gorm.DB
}

// InitConnectionMysql open initialize a new db connection.
func InitConnectionMysql(cf *config.Configs) (*MysqlConfig, error) {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True",
		cf.Mysql.Username,
		cf.Mysql.Password,
		cf.Mysql.Host,
		cf.Mysql.Port,
		cf.Mysql.DatabaseName,
	)

	database, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		logrus.Errorf("[InitConnectionMysql] failed to connect to the database error: %s", err)
		return nil, err
	}

	sqlDB, err := database.DB()
	if err != nil {
		logrus.Errorf("[InitConnectionMysql] set up to connect to the database error: %s", err)
		return nil, err
	}

	sqlDB.SetMaxOpenConns(cf.Mysql.ConnectionPool.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cf.Mysql.ConnectionPool.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cf.Mysql.ConnectionPool.MaxLifeTime * time.Minute)

	return &MysqlConfig{Database: database}, nil
}

// ExportDatabase export database to global variable.
func (config *MysqlConfig) ExportDatabase() {
	Database = config.Database
}
