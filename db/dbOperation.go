package db

import (
	"database/sql"
	"github.com/BlackofZero/go-models/errors"
)

type ExecInstance interface {
	Query(database, statement string) (*sql.Rows, errors.Error)
	ParseRows(rows *sql.Rows) ([]string, [][]string, errors.Error)
	BatchExec(database string, statements []string) errors.Error
	GetMinMax(min string, max string) (string, string, bool, errors.Error)
}
