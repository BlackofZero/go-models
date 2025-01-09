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

func (m mysqlExec) connect(database string) (*gorm.DB, *sql.DB, errors.Error) {
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

	if m.sqldb == nil {
		db, err := gorm.Open(drivermysql.New(initMySQLConfig(dsn)))
		if err != nil {
			return nil, nil, errors.New(err.Error())
		}
		m.sqldb, err = db.DB()
		if err != nil {
			return nil, nil, errors.New(err.Error())
		}
		m.sqldb.SetMaxIdleConns(5)
		m.sqldb.SetMaxOpenConns(20)
		m.sqldb.SetConnMaxLifetime(time.Minute * 3)
		return db, m.sqldb, nil
	}
	return nil, m.sqldb, nil
}

func (m mysqlExec) Connect(database string) (*gorm.DB, *sql.DB, errors.Error) {
	return m.connect(database)
}

func (m mysqlExec) QueryContext(ctx context.Context, database, statement string) (*sql.Rows, errors.Error) {
	_, db, err := m.Connect(database)
	if err != nil {
		return nil, err
	} else {
		defer db.Close()
		rows, err := db.QueryContext(ctx, statement)
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
	_, db, err := m.Connect(database)
	if err != nil {
		return err
	} else {
		defer db.Close()
		for _, s := range statements {
			_, err := db.ExecContext(ctx, s)
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
}

func (m mysqlExec) BatchExec(database string, statements []string) errors.Error {
	return m.batchExec(context.Background(), database, statements)
}
