package db

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host string
	Port string
	Login string
	Pass string
	Db string
}

var Connect *sqlx.DB

func ConnectMySQL(cfg Config) *sqlx.DB{
	connectString := []string{cfg.Login, ":", cfg.Pass, "@tcp(", cfg.Host, cfg.Port, ")/", cfg.Db}
	connect, err := sqlx.Connect("mysql", strings.Join(connectString, ""))
	if err != nil {
		panic(err)
	}

	return connect
}