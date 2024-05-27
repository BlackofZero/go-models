package mysql

import "github.com/BlackofZero/go-models/db"

type mysqlExec struct {
	hostIP   string
	port     int
	username string
	password string
	db       string
	table    table
}
type table struct {
	tablename  string
	condition  string
	primarykey string
	limit      int
}

type Mysql struct {
	Url       string    `json:"url"`
	Port      int       `json:"port"`
	Db        string    `json:"db"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Operation Operation `json:"operation"`
}
type Operation struct {
	Table      string `json:"table"`
	Sql        string `json:"sql"`
	Limit      int    `json:"limit"`
	PrimaryKey string `json:"primaryKey"`
	Start      string `json:"start"`
	End        string `json:"end"`
	Condition  string `json:"condition"`
}

func NewMysqlExec(mysql Mysql) db.ExecInstance {
	return &mysqlExec{hostIP: mysql.Url, port: mysql.Port, username: mysql.Username, password: mysql.Password, db: mysql.Db,
		table: table{tablename: mysql.Operation.Table, condition: mysql.Operation.Condition, limit: mysql.Operation.Limit, primarykey: mysql.Operation.PrimaryKey}}
}
