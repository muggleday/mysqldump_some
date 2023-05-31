package mysqldump_some

import (
	"database/sql"
	"time"
)

// TABLE mapped from table <TABLES>
type TABLE struct {
	TABLECATALOG   *string    `gorm:"column:TABLE_CATALOG" json:"TABLE_CATALOG"`
	TABLESCHEMA    *string    `gorm:"column:TABLE_SCHEMA" json:"TABLE_SCHEMA"`
	TABLENAME      *string    `gorm:"column:TABLE_NAME" json:"TABLE_NAME"`
	TABLETYPE      string     `gorm:"column:TABLE_TYPE;not null" json:"TABLE_TYPE"`
	ENGINE         *string    `gorm:"column:ENGINE" json:"ENGINE"`
	VERSION        *int32     `gorm:"column:VERSION" json:"VERSION"`
	ROWFORMAT      *string    `gorm:"column:ROW_FORMAT" json:"ROW_FORMAT"`
	TABLEROWS      *int64     `gorm:"column:TABLE_ROWS" json:"TABLE_ROWS"`
	AVGROWLENGTH   *int64     `gorm:"column:AVG_ROW_LENGTH" json:"AVG_ROW_LENGTH"`
	DATALENGTH     *int64     `gorm:"column:DATA_LENGTH" json:"DATA_LENGTH"`
	MAXDATALENGTH  *int64     `gorm:"column:MAX_DATA_LENGTH" json:"MAX_DATA_LENGTH"`
	INDEXLENGTH    *int64     `gorm:"column:INDEX_LENGTH" json:"INDEX_LENGTH"`
	DATAFREE       *int64     `gorm:"column:DATA_FREE" json:"DATA_FREE"`
	AUTOINCREMENT  *int64     `gorm:"column:AUTO_INCREMENT" json:"AUTO_INCREMENT"`
	CREATETIME     time.Time  `gorm:"column:CREATE_TIME;not null" json:"CREATE_TIME"`
	UPDATETIME     *time.Time `gorm:"column:UPDATE_TIME" json:"UPDATE_TIME"`
	CHECKTIME      *time.Time `gorm:"column:CHECK_TIME" json:"CHECK_TIME"`
	TABLECOLLATION *string    `gorm:"column:TABLE_COLLATION" json:"TABLE_COLLATION"`
	CHECKSUM       *int64     `gorm:"column:CHECKSUM" json:"CHECKSUM"`
	CREATEOPTIONS  *string    `gorm:"column:CREATE_OPTIONS" json:"CREATE_OPTIONS"`
	TABLECOMMENT   *string    `gorm:"column:TABLE_COMMENT" json:"TABLE_COMMENT"`
}

func QueryTables(db *sql.DB, database string) ([]TABLE, error) {
	rows, err := db.Query(`
		SELECT TABLE_NAME
		FROM information_schema.TABLES
		WHERE TABLE_SCHEMA = ?;
		`,
		database,
	)
	if err != nil {
		return nil, err
	}

	var got []TABLE
	for rows.Next() {
		var r TABLE
		err = rows.Scan(&r.TABLENAME)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}
