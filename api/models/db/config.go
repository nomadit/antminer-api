package db

import (
	"fmt"
	"github.com/nomadit/antminer-api/config"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const (
	RunStatus            = "RUN"
	RebootStatus         = "REBOOT"
	ErrorNoWorker        = "ERROR_NO_WORKER"
	ErrorHashRate        = "ERROR_HASH_RATE"
	ErrorOverTemperature = "ERROR_OVER_TEMPERATURE"
	CommandChangeConfig  = "CHANG_CONFIGURE"
)

var (
	conf *Config
	once sync.Once
)

// Config is instance for aws rds access
type Config struct {
	rdb *sqlx.DB
}

// GetAWSDB return sql.DB
func NewDB(db *config.DBConfig) *sqlx.DB {
	once.Do(func() {
		conf = new(Config)
		conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", db.User, db.Password, db.Host, db.DB)
		var err error
		conf.rdb, err = sqlx.Open("mysql", conn)
		if err != nil {
			panic(err)
		}
	})
	return conf.rdb.Unsafe()
}

func GetDB() *sqlx.DB {
	return conf.rdb.Unsafe()
}
