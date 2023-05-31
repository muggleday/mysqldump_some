package mysqldump_some

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

type Dumper struct {
	dsn string

	db *sql.DB

	Database    string
	OutDir      string
	DumpOptions []string

	ExcludedTables []string
	IncludedTables []string

	TableConditions    []DumpTableCondition
	MatchAnyConditions []DumpMatchCondition // match any
}

func NewDumper(dsn, database string) (*Dumper, error) {
	dumper := &Dumper{
		dsn:      dsn,
		Database: database,
		OutDir:   fmt.Sprintf("%s_%s", database, time.Now().Format("20060102150405")),
	}

	var err error
	dumper.db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = dumper.db.Ping()
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(dumper.OutDir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return dumper, nil
}
