package mysql

import (
	"fmt"
	"github.com/BlackofZero/go-models/db"
	"sync"
)

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

var rwLock sync.RWMutex
var exportRows int

func SetExportRows(rows int) {
	rwLock.Lock()
	exportRows = rows
	rwLock.Unlock()
}

func GetExportRows() int {
	rwLock.RLock()
	rows := exportRows
	rwLock.RUnlock()
	return rows
}

func NewMysqlExec(mysql Mysql) db.ExecInstance {
	return &mysqlExec{hostIP: mysql.Url, port: mysql.Port, username: mysql.Username, password: mysql.Password, db: mysql.Db,
		table: table{tablename: mysql.Operation.Table, condition: mysql.Operation.Condition, limit: mysql.Operation.Limit, primarykey: mysql.Operation.PrimaryKey}}
}

func NewExec(mysql Mysql) db.ExecInstance {
	return &mysqlExec{hostIP: mysql.Url, port: mysql.Port, username: mysql.Username, password: mysql.Password, db: mysql.Db}
}

func AssembleSql(operation Operation, min, max string) string {
	if min == max {
		return fmt.Sprintf(
			"%s WHERE `%s` = '%s' AND %s ORDER BY `%s` LIMIT %d;",
			operation.Sql,
			operation.PrimaryKey,
			max,
			operation.Condition,
			operation.PrimaryKey,
			operation.Limit,
		)
	}

	return fmt.Sprintf(
		"%s WHERE `%s` >= '%s' AND `%s` < '%s' AND %s ORDER BY `%s` LIMIT %d;",
		operation.Sql,
		operation.PrimaryKey,
		min,
		operation.PrimaryKey,
		max,
		operation.Condition,
		operation.PrimaryKey,
		operation.Limit,
	)
}
