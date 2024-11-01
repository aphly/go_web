package connect

import (
	"fmt"
	"go_web/app/core/config"
	"go_web/app/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

var dbs = make(map[string]*gorm.DB)
var dbLocker sync.RWMutex

func Mysql(config *config.Db, key string) *gorm.DB {
	dbLocker.RLock()
	db, ok := dbs[key]
	if ok {
		dbLocker.RUnlock()
		return db
	}
	dbLocker.RUnlock()

	dbLocker.Lock()
	defer dbLocker.Unlock()
	if _, ok1 := dbs[key]; ok1 {
		return dbs[key]
	}
	dbs[key] = getDb(config)
	return dbs[key]
}

func getDb(config *config.Db) *gorm.DB {
	var LWriter helper.WriterLog
	newLogger := logger.New(
		LWriter,
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local&timeout=%dms&writeTimeout=%dms&readTimeout=%dms",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.Charset,
		config.TimeOut,
		config.WriteTimeOut,
		config.ReadTimeOut,
	)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:               dsn,
		DefaultStringSize: 255,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "",   // 表前缀
			SingularTable: true, // 禁用表名复数
		},
		Logger: newLogger,
	})
	if err != nil {
		panic("连接数据库失败, error:" + err.Error())
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic("数据库池, error:" + err.Error())
	}
	sqlDB.SetMaxIdleConns(config.MaxIdleConnect)
	sqlDB.SetMaxOpenConns(config.MaxOpenConnect)
	sqlDB.SetConnMaxLifetime(time.Second * 120)
	return db
}
