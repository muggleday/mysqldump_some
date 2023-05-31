package mysqldump_some

import "database/sql"

// COLUMN mapped from table <COLUMNS>
type COLUMN struct {
	TABLECATALOG           *string `gorm:"column:TABLE_CATALOG" json:"TABLE_CATALOG"`
	TABLESCHEMA            *string `gorm:"column:TABLE_SCHEMA" json:"TABLE_SCHEMA"`
	TABLENAME              *string `gorm:"column:TABLE_NAME" json:"TABLE_NAME"`
	COLUMNNAME             *string `gorm:"column:COLUMN_NAME" json:"COLUMN_NAME"`
	ORDINALPOSITION        int32   `gorm:"column:ORDINAL_POSITION;not null" json:"ORDINAL_POSITION"`
	COLUMNDEFAULT          *string `gorm:"column:COLUMN_DEFAULT" json:"COLUMN_DEFAULT"`
	ISNULLABLE             string  `gorm:"column:IS_NULLABLE;not null" json:"IS_NULLABLE"`
	DATATYPE               *string `gorm:"column:DATA_TYPE" json:"DATA_TYPE"`
	CHARACTERMAXIMUMLENGTH *int64  `gorm:"column:CHARACTER_MAXIMUM_LENGTH" json:"CHARACTER_MAXIMUM_LENGTH"`
	CHARACTEROCTETLENGTH   *int64  `gorm:"column:CHARACTER_OCTET_LENGTH" json:"CHARACTER_OCTET_LENGTH"`
	NUMERICPRECISION       *int64  `gorm:"column:NUMERIC_PRECISION" json:"NUMERIC_PRECISION"`
	NUMERICSCALE           *int64  `gorm:"column:NUMERIC_SCALE" json:"NUMERIC_SCALE"`
	DATETIMEPRECISION      *int32  `gorm:"column:DATETIME_PRECISION" json:"DATETIME_PRECISION"`
	CHARACTERSETNAME       *string `gorm:"column:CHARACTER_SET_NAME" json:"CHARACTER_SET_NAME"`
	COLLATIONNAME          *string `gorm:"column:COLLATION_NAME" json:"COLLATION_NAME"`
	COLUMNTYPE             string  `gorm:"column:COLUMN_TYPE;not null" json:"COLUMN_TYPE"`
	COLUMNKEY              string  `gorm:"column:COLUMN_KEY;not null" json:"COLUMN_KEY"`
	EXTRA                  *string `gorm:"column:EXTRA" json:"EXTRA"`
	PRIVILEGES             *string `gorm:"column:PRIVILEGES" json:"PRIVILEGES"`
	COLUMNCOMMENT          string  `gorm:"column:COLUMN_COMMENT;not null" json:"COLUMN_COMMENT"`
	GENERATIONEXPRESSION   string  `gorm:"column:GENERATION_EXPRESSION;not null" json:"GENERATION_EXPRESSION"`
	SRSID                  *int32  `gorm:"column:SRS_ID" json:"SRS_ID"`
}

func QueryTableColumns(db *sql.DB, database, tableName string) ([]COLUMN, error) {
	rows, err := db.Query(`
		SELECT COLUMN_NAME
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? 
		  AND TABLE_NAME = ?;
		`,
		database,
		tableName,
	)
	if err != nil {
		return nil, err
	}

	var got []COLUMN
	for rows.Next() {
		var r COLUMN
		err = rows.Scan(&r.COLUMNNAME)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}

func IsColumnExistsInTable(db *sql.DB, database, tableName, columnName string) (bool, error) {
	existFlag := false
	err := db.QueryRow(`
		SELECT CASE COUNT(1) WHEN 0 THEN FALSE ELSE TRUE END
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = ? 
		  AND TABLE_NAME = ? 
		  AND COLUMN_NAME = ? 
		LIMIT 1;
		`,
		database,
		tableName,
		columnName,
	).Scan(&existFlag)
	if err != nil {
		return false, err
	}

	return existFlag, nil
}
