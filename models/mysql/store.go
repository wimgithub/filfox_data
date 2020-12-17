package mysql

import (
	"filfox_data/models"
	"filfox_data/pkg/logging"
	"filfox_data/pkg/setting"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"reflect"
	"sync"
	"time"
)

var gdb *gorm.DB
var store model.Store
var storeOnce sync.Once

type Store struct {
	db *gorm.DB
}

func SharedStore() model.Store {
	storeOnce.Do(func() {
		err := initDb()
		if err != nil {
			panic(err)
		}
		store = NewStore(gdb)
	})
	return store
}

func NewStore(db *gorm.DB) *Store {
	return &Store{
		db: db,
	}
}

func initDb() (err error) {
	gdb, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=5s&readTimeout=3s&writeTimeout=3s&parseTime=true&loc=Local&charset=utf8mb4,utf8",
		setting.DatabaseSetting.User,
		setting.DatabaseSetting.Password,
		setting.DatabaseSetting.Host,
		setting.DatabaseSetting.Name))
	if err != nil {
		return err
	}
	// 创建数据库的时候名字不是复数
	gdb.SingularTable(true)
	// SetMaxIdleCons 设置连接池中的最大闲置连接数。
	gdb.DB().SetMaxIdleConns(10)
	// SetMaxOpenCons 设置数据库的最大连接数量。
	gdb.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetiment 设置连接的最大可复用时间。
	gdb.DB().SetConnMaxLifetime(time.Hour)

	// 配置表名前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return setting.DatabaseSetting.TablePrefix + defaultTableName
	}

	var tables = []interface{}{
		&model.Data{},
	}
	for _, table := range tables {
		logging.Info("migrating database, table: %v", reflect.TypeOf(table))
		if err = gdb.AutoMigrate(table).Error; err != nil {
			return err
		}
	}

	return
}

func (s *Store) BeginTx() (model.Store, error) {
	db := s.db.Begin()
	if db.Error != nil {
		return nil, db.Error
	}
	return NewStore(db), nil
}

func (s *Store) Rollback() error {
	return s.db.Rollback().Error
}

func (s *Store) CommitTx() error {
	return s.db.Commit().Error
}
