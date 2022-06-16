package mysql

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/sioncojp/famili-api/utils/config"
	"github.com/sioncojp/famili-api/utils/log"
)

// NewMySQLHandler...MySQLとコネクションする
func NewMySQLHandler(c *config.DataStoreConfig) (*gorm.DB, error) {
	log.Log.Debug("new infrastructure MySQLHandler")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username,
		c.Password,
		c.Url,
		c.Port,
		c.DbName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
