package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/BlackofZero/go-models/errors"
	drivermysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func (m mysqlExec) connect(database string) errors.Error {
	m.mutex.Lock() // 加锁，确保线程安全
	defer m.mutex.Unlock()
	if Sqldb != nil {
		return nil
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=10s&allowAllFiles=true",
		m.username, m.password, m.hostIP, m.port, database,
	)
	//filename := "mysql.log"
	//logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	//if err != nil {
	//	return nil, nil, errors.New(err.Error())
	//}
	//db, err := gorm.Open(drivermysql.New(initMySQLConfig(dsn)), initGormConfig(logger.Info, logFile))

	db, err := gorm.Open(drivermysql.New(initMySQLConfig(dsn)))
	if err != nil {
		return errors.New(err.Error())
	}
	sqldb, err := db.DB()
	if err != nil {
		return errors.New(err.Error())
	}
	sqldb.SetMaxIdleConns(5)
	sqldb.SetMaxOpenConns(20)
	sqldb.SetConnMaxLifetime(time.Minute * 1)
	Sqldb = sqldb
	return nil

}

// 提供关闭连接的方法
func (m *mysqlExec) Close() errors.Error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if Sqldb != nil {
		err := Sqldb.Close()
		Sqldb = nil
		return errors.New(err.Error())
	}
	return nil
}
func (m mysqlExec) Connect(database string) errors.Error {
	return m.connect(database)
}

func (m mysqlExec) QueryContext(ctx context.Context, database, statement string) (*sql.Rows, errors.Error) {
	if Sqldb == nil {
		m.Connect(database)
	}
	rows, err := Sqldb.QueryContext(ctx, statement)
	if err != nil {
		es := errors.New(fmt.Sprintf("执行SQL语句报错: [%s], %s", statement, err.Error()))
		if err.Error() == context.Canceled.Error() {
			es = errors.New("SQL执行终止")
		} else if err.Error() == context.DeadlineExceeded.Error() {
			es = errors.New("SQL执行超时")
		}
		return rows, es
	}
	return rows, nil

}
func (m mysqlExec) GetMinMax(min, max string) (string, string, bool, errors.Error) {
	overFalg := false

	whereCondition := m.table.condition
	if "" != min {
		whereCondition += fmt.Sprintf(" AND `%s` >= '%s'", m.table.primarykey, min)
	}
	if "" != max {
		whereCondition += fmt.Sprintf(" AND `%s` < '%s'", m.table.primarykey, max)
	}
	rows, err := m.Query(
		m.db,
		fmt.Sprintf(
			"SELECT min(`%s`), max(`%s`) FROM (SELECT `%s` FROM `%s` WHERE %s ORDER BY `%s` LIMIT %d) a;",
			m.table.primarykey,
			m.table.primarykey,
			m.table.primarykey,
			m.table.tablename,
			whereCondition,
			m.table.primarykey,
			m.table.limit,
		),
	)
	if err != nil {
		return "", "", overFalg, err
	}
	_, results, err := m.ParseRows(rows)
	if err != nil {
		return "", "", overFalg, err
	}
	// no return use End as max
	//if len(results[0]) < m.table.limit-1 {
	//	overFalg = true
	//}

	if results[0][0] == results[0][1] {
		overFalg = true

	}
	if len(results) == 0 {
		overFalg = true
		return "", "", overFalg, nil

	}
	return results[0][0], results[0][1], overFalg, nil
}

func (m mysqlExec) Query(database, statement string) (*sql.Rows, errors.Error) {
	return m.QueryContext(context.Background(), database, statement)
}

func (m mysqlExec) ParseRows(rows *sql.Rows) ([]string, [][]string, errors.Error) {
	columns, err := rows.Columns()
	if err != nil {
		return []string{}, [][]string{}, errors.New("解析返回结果报错: " + err.Error())
	}
	var results [][]string
	for rows.Next() {
		var r []interface{}
		for i := 0; i < len(columns); i++ {
			r = append(r, &[]byte{})
		}
		err := rows.Scan(r...)
		if err != nil {
			return []string{}, [][]string{}, errors.New("解析返回结果报错: " + err.Error())

		}
		var result []string
		for _, v := range r {
			result = append(result, string(*v.(*[]byte)))
		}
		results = append(results, result)
	}
	return columns, results, nil
}

func (m mysqlExec) batchExec(ctx context.Context, database string, statements []string) errors.Error {
	if Sqldb == nil {
		m.Connect(database)
	}
	for _, s := range statements {
		_, err := Sqldb.ExecContext(ctx, s)
		if err != nil {
			es := errors.New(fmt.Sprintf("执行SQL语句报错: [%s], %s", statements, err.Error()))
			if err.Error() == context.Canceled.Error() {
				es = errors.New("SQL执行终止")
			} else if err.Error() == context.DeadlineExceeded.Error() {
				es = errors.New("SQL执行超时")
			}
			return es
		}
	}
	return nil

}

func (m mysqlExec) BatchExec(database string, statements []string) errors.Error {
	return m.batchExec(context.Background(), database, statements)
}
