package x

import (
	"fmt"
	"net/url"

	"github.com/spf13/viper"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 根据平台XML数据库连接:x.Conndb()
// Global database use gorm
var Gdb *gorm.DB

var err error

func Conndb() {
	var dbv *viper.Viper
	if dbv, err = DatabaseXmlToViper(); err != nil {
		Xlog.Error(err.Error())
		return
	}
	if Gdb, err = ConnGormMssql(dbv); err != nil {
		Xlog.Error(err.Error())
		return
	}

	fmt.Println("db connect success")
}

// SetDefaultViperConfig sets the default configuration for viper
func SetDefaultViperConfig(v *viper.Viper) {
	v.SetDefault("usr", "defaultUser")
	v.SetDefault("pwd", "defaultPassword")
	v.SetDefault("host", "localhost")
	v.SetDefault("port", 1433)
	v.SetDefault("dbname", "defaultDB")
	v.SetDefault("dbextra", "encrypt=true")
	v.SetDefault("maxIdleConns", 10)
	v.SetDefault("maxOpenConns", 100)
}

// Connect to the mssqlserver database, and the caller decides how to handle the error
func ConnGormMssql(initCfg *viper.Viper) (*gorm.DB, error) {
	// URL-encode passwords
	encodedPwd := url.QueryEscape(initCfg.GetString("pwd"))
	// sqlserver://gorm:LoremIpsum86@localhost:9930?database=gorm
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&%s",
		initCfg.GetString("usr"), encodedPwd, initCfg.GetString("host"),
		initCfg.GetInt("port"), initCfg.GetString("dbname"),
		initCfg.GetString("dbextra"))
	mssqlConfig := sqlserver.Config{
		DSN: dsn, // DSN data source name
	}
	if db, err := gorm.Open(sqlserver.New(mssqlConfig), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "T_", // table name prefix, table for `User` would be `t_users`
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		DisableForeignKeyConstraintWhenMigrating: true, //注意 AutoMigrate 会自动创建数据库外键约束，您可以在初始化时禁用此功能
	}); err != nil {
		return nil, err
	} else {
		// 其他初始化

		// sqlDB, _ := db.DB()
		// sqlDB.SetMaxIdleConns(initCfg.GetInt("maxIdleConns"))
		// sqlDB.SetMaxOpenConns(initCfg.GetInt("maxOpenConns"))
		// fmt.Println("init in person subsidiary")
		// err := db.Debug().AutoMigrate(&PersonNewProperty{})
		// if err != nil {
		// 	Xlog.Error("Gdb.Debug().AutoMigrate() Error", zap.Error(err))
		// }

		return db, nil
	}
}

// type PersonNewProperty struct {
// 	Organization        string `gorm:"column:organization"`
// 	OrganizationId      int64  `gorm:"column:organization_id"`
// 	OrganizationClassId int64  `gorm:"column:organization_classid"`
// 	Subsidiary          string `gorm:"column:subsidiary"`
// 	SubsidiaryObjectId  int64  `gorm:"column:subsidiary_objectid"`
// 	SubsidiaryClassId   int64  `gorm:"column:subsidiary_classid"`
// }

// func (*PersonNewProperty) TableName() string {
// 	return "T_Person"
// }
